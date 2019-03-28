package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/coast-team/mute-auth-proxy/helper"
	"github.com/dgraph-io/badger"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// UserAllPK is the structure that contains all the public Keys of an user (one PK per device)
type UserAllPK struct {
	Login string
	AllPK map[string]string
}

// PublicKey represents a public key in a JSON object
type PublicKey struct {
	PK string `json:"pk"`
}

// UserPublicKey is the structure that contains the public key associate to an user and a device
type UserPublicKey struct {
	Login  string `json:"login"`
	Device string `json:"deviceID"`
	PK     string `json:"pk"`
}

type pkAlreadyExistsError struct {
	login  string
	device string
}

func (e *pkAlreadyExistsError) Error() string {
	return fmt.Sprintf("PK already present for %s-%s", e.login, e.device)
}

type deviceAlreadyExistsError struct {
	login string
}

func (e *deviceAlreadyExistsError) Error() string {
	return fmt.Sprintf("Device already present for %s", e.login)
}

// MakePublicKeyPOSTHandler is the handler for the API to save a public key
// This public key is associated to an username and a deviceID
func MakePublicKeyPOSTHandler(db *badger.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "Please send a request body", http.StatusBadRequest)
			return
		}
		var userPK UserPublicKey
		err := json.NewDecoder(r.Body).Decode(&userPK)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Printf("Keyserver ADD, error while parsing JSON: %s\n", err)
			return
		}
		err = validateJWT(r, userPK.Login, true)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Printf("Keyserver ADD, JWT validation err: %s\n", err)
			return
		}
		addErr := handleAddPublicKey(db, userPK.Login, userPK.Device, userPK.PK)
		if addErr != nil {
			if err, ok := addErr.(*pkAlreadyExistsError); ok {
				http.Error(w,
					fmt.Sprintf("%s - PK already registered for %s:%s", http.StatusText(http.StatusBadRequest), userPK.Login, userPK.Device),
					http.StatusBadRequest)
				log.Printf("Keyserver ADD err (pk already exists): %s\n", err)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Printf("Keyserver ADD err: %s\n", addErr)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", fmt.Sprintf("/public-key/%s/%s", userPK.Login, userPK.Device))
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(userPK)
		if err != nil {
			log.Printf("Keyserver ADD err (response marshalling): %s\n", err)
		}
	}
}

func handleAddPublicKey(db *badger.DB, login, device, pk string) error {
	log.Printf("Wanting to add PK %s for %s:%s", pk, login, device)
	found, err := checkPKEntryAlreadyExists(db, login, device)
	if err != nil {
		return err
	}
	if found {
		return &pkAlreadyExistsError{login, device}
	}
	found, err = checkPKListEntryAlreadyExists(db, login, device)
	if err != nil {
		return err
	}
	if found {
		return &deviceAlreadyExistsError{login}
	}
	err = db.View(makeBDAddTxnHandler(db, login, device, pk))
	return err
}

func checkPKEntryAlreadyExists(db *badger.DB, login, device string) (b bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(fmt.Sprintf("%s:%s", login, device)))
		if err != nil {
			return err
		}
		pk, copyErr := item.ValueCopy(nil)
		if copyErr != nil {
			panic(copyErr)
		}
		log.Printf("Check PK entry : %s:%s - %s", login, device, pk)
		return nil
	})
	if err == badger.ErrKeyNotFound {
		return false, nil
	}
	return true, err
}

func checkPKListEntryAlreadyExists(db *badger.DB, login, device string) (b bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	var devices []string
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(login))
		if err != nil {
			return err
		}
		pkList, copyErr := item.ValueCopy(nil)
		if copyErr != nil {
			panic(copyErr)
		}
		jsonErr := json.Unmarshal(pkList, &devices)
		if jsonErr != nil {
			panic(jsonErr)
		}
		log.Printf("Check devices entry : %s - %s", login, devices)
		return nil
	})
	if err == badger.ErrKeyNotFound {
		return false, nil
	}
	if !helper.StringInSlice(device, devices) {
		return false, nil
	}
	return true, nil
}

func makeBDAddTxnHandler(db *badger.DB, login, device, pk string) func(txn *badger.Txn) error {
	return func(txn *badger.Txn) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()
		var devices []string
		getErr := db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(login))
			if err != nil {
				return err
			}
			pkList, copyErr := item.ValueCopy(nil)
			if copyErr != nil {
				panic(copyErr)
			}
			jsonErr := json.Unmarshal(pkList, &devices)
			if jsonErr != nil {
				panic(jsonErr)
			}
			return nil
		})
		if getErr != nil && getErr != badger.ErrKeyNotFound {
			return getErr
		}
		devices = append(devices, device)
		db.Update(func(txn *badger.Txn) error {
			setErr := txn.Set([]byte(fmt.Sprintf("%s:%s", login, device)), []byte(pk))
			if setErr != nil {
				panic(setErr)
			}
			value, jsonErr := json.Marshal(devices)
			if jsonErr != nil {
				panic(jsonErr)
			}
			setErr = txn.Set([]byte(fmt.Sprintf("%s", login)), value)
			if setErr != nil {
				panic(setErr)
			}
			return nil
		})
		return
	}
}

// MakePublicKeyGETHandler is the handler for the API to get a public key from a specific user and deviceID
func MakePublicKeyGETHandler(db *badger.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		login := vars["login"]
		device := vars["device"]
		err := validateJWT(r, login, false)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Printf("Keyserver GET, JWT validation err: %s\n", err)
			return
		}
		pk, err := handleGetPublicKey(db, login, device)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				http.Error(w,
					fmt.Sprintf("%s - PK not found for %s:%s", http.StatusText(http.StatusNotFound), login, device),
					http.StatusNotFound)
				log.Printf("Keyserver GET err : PK not found for %s:%s\n", login, device)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Printf("Keyserver GET err: %s\n", err)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(pk)
		if err != nil {
			log.Printf("Keyserver GET err (response marshalling): %s\n", err)
		}
	}
}

func handleGetPublicKey(db *badger.DB, login, device string) (PublicKey, error) {
	log.Printf("Get PK for : %s-%s", login, device)
	var pk PublicKey
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(fmt.Sprintf("%s:%s", login, device)))
		if err != nil {
			return err
		}
		dbPK, copyErr := item.ValueCopy(nil)
		pk.PK = string(dbPK)
		if copyErr != nil {
			return copyErr
		}
		log.Printf("Check PK entry : %s:%s - %s", login, device, pk)
		return nil
	})
	return pk, err
}

// MakePublicKeyGETAllHandler is the handler for the API to get all the public keys of an user
func MakePublicKeyGETAllHandler(db *badger.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		login := vars["login"]
		device := vars["device"]
		err := validateJWT(r, login, true)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Printf("Keyserver GET ALL, JWT validation err: %s\n", err)
			return
		}
		allPK, err := handleGetAllPublicKeys(db, login)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				http.Error(w,
					fmt.Sprintf("%s - PK not found for %s:%s", http.StatusText(http.StatusNotFound), login, device),
					http.StatusNotFound)
				log.Printf("Keyserver GET ALL err : PK not found for %s:%s\n", login, device)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Printf("Keyserver GET ALL err: %s\n", err)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(allPK)
		if err != nil {
			log.Printf("Keyserver GET err (response marshalling): %s\n", err)
		}
	}
}

func handleGetAllPublicKeys(db *badger.DB, login string) (UserAllPK, error) {
	log.Printf("Get devices for : %s", login)
	var userAllPK UserAllPK
	userAllPK.Login = login
	userAllPK.AllPK = make(map[string]string)
	err := db.View(func(txn *badger.Txn) error {
		var deviceList []string
		item, err := txn.Get([]byte(fmt.Sprintf("%s", login)))
		if err != nil {
			return err
		}
		devices, copyErr := item.ValueCopy(nil)
		if copyErr != nil {
			return copyErr
		}
		jsonErr := json.Unmarshal(devices, &deviceList)
		if jsonErr != nil {
			return jsonErr
		}
		log.Printf("Check devices entry : %s - %s", login, deviceList)
		for _, device := range deviceList {
			item, err := txn.Get([]byte(fmt.Sprintf("%s:%s", login, device)))
			if err != nil {
				return err
			}
			pk, copyErr := item.ValueCopy(nil)
			if copyErr != nil {
				return copyErr
			}
			log.Printf("Check PK entry : %s:%s - %s", login, device, pk)
			userAllPK.AllPK[device] = string(pk)
		}
		return nil
	})
	log.Printf("Check All PK for %s - %#v", userAllPK.Login, userAllPK.AllPK)
	return userAllPK, err
}

// MakePublicKeyPUTHandler is the handler for the API to update a public key from an username deviceID
func MakePublicKeyPUTHandler(db *badger.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		login := vars["login"]
		device := vars["device"]
		err := validateJWT(r, login, true)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Printf("Keyserver(PUT) JWT validation err: %s\n", err)
			return
		}
		var pk PublicKey
		err = json.NewDecoder(r.Body).Decode(&pk)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Printf("Keyserver UPDATE, error while parsing JSON: %s\n", err)
			return
		}
		err = handleUpdatePublicKeys(db, login, device, pk.PK)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				http.Error(w,
					fmt.Sprintf("%s - PK not found for %s:%s", http.StatusText(http.StatusNotFound), login, device),
					http.StatusNotFound)
				log.Printf("Keyserver UPDATE err : PK not found for %s:%s\n", login, device)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Printf("Keyserver update err: %s\n", err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleUpdatePublicKeys(db *badger.DB, login, device, pk string) error {
	log.Printf("Update PK for %s:%s, new value: %s", login, device, pk)
	found, err := checkPKEntryAlreadyExists(db, login, device)
	if err != nil {
		return err
	}
	if found {
		err = db.Update(func(txn *badger.Txn) error {
			setErr := txn.Set([]byte(fmt.Sprintf("%s:%s", login, device)), []byte(pk))
			if setErr != nil {
				panic(setErr)
			}
			return nil
		})
	} else {
		err = badger.ErrKeyNotFound
	}
	return err
}

func makeBDUpdateTxnHandler(db *badger.DB, login, device, pk string) func(txn *badger.Txn) error {
	return func(txn *badger.Txn) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()
		var devices []string
		getErr := db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(login))
			if err != nil {
				return err
			}
			pkList, copyErr := item.ValueCopy(nil)
			if copyErr != nil {
				panic(copyErr)
			}
			jsonErr := json.Unmarshal(pkList, &devices)
			if jsonErr != nil {
				panic(jsonErr)
			}
			return nil
		})
		if getErr != nil && getErr != badger.ErrKeyNotFound {
			return getErr
		}
		devices = append(devices, device)
		db.Update(func(txn *badger.Txn) error {
			setErr := txn.Set([]byte(fmt.Sprintf("%s:%s", login, device)), []byte(pk))
			if setErr != nil {
				panic(setErr)
			}
			value, jsonErr := json.Marshal(devices)
			if jsonErr != nil {
				panic(jsonErr)
			}
			setErr = txn.Set([]byte(fmt.Sprintf("%s", login)), value)
			if setErr != nil {
				panic(setErr)
			}
			return nil
		})
		return
	}
}

func validateJWT(r *http.Request, login string, checkLogin bool) error {
	token, err := helper.ExtractJWT(r)
	if err != nil {
		err = helper.IsJWTValid(token, err)
		return fmt.Errorf("Couldn't extract or validate the JWT.\nError was: %s", err)
	}
	if checkLogin {
		err = validateLogin(login, token.Claims.(jwt.MapClaims)["login"].(string))
		if err != nil {
			return fmt.Errorf("Unallowed get, creation or modification of a PK.\nError was: %s", err)
		}
	}
	return nil
}

func validateLogin(login string, tokenLogin string) error {
	if login == tokenLogin {
		return nil
	} else if login == fmt.Sprintf("%s@github", tokenLogin) {
		return nil
	}

	return fmt.Errorf("Difference between connected login and login in the API request\nLogin : %s\nLogin in request : %s", tokenLogin, login)
}
