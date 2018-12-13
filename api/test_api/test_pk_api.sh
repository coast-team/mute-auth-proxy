#! /usr/bin/env sh

set -e


if [ $# -ne 1 ]; then
    echo "Usage : ./test_pk_api <jwt>"
    exit 1
fi

mkdir -p /tmp/test_pk_api/

baseUrl="http://localhost:4000"


echo -e "=========== Begin Test 1 ===========\nTry to add the PK of someone else than you (given the JWT that identifies yourself)\n"
loginWrong="Alice"
deviceWrong="device1"
pkWrong="pk1"
curl "$baseUrl/public-key" -d '{"login":"'"$loginWrong"'","device":"'"$deviceWrong"'","pk":"'"$pkWrong"'"}' -H "Content-Type: application/json" -H "authorization: Bearer $1" -i
echo -e "\n=========== End Test 1 ===========\nResponse should be HTTP error 401 Unauthorized\n\n"

loginJP="jpeisenbarth"
deviceJP1="device1"
deviceJP2="device2"
pkJP1="pk1"
pkJP2="pk2"
echo -e "=========== Begin Test 2 ===========\nTry to add PK ($pkJP1, $pkJP2) for $loginJP-$deviceJP1 and $loginJP-$deviceJP2\n"
curl "$baseUrl/public-key" -d '{"login":"'"$loginJP"'","device":"'"$deviceJP1"'","pk":"'"$pkJP1"'"}' -H "Content-Type: application/json" -H "authorization: Bearer $1" -i
curl "$baseUrl/public-key" -d '{"login":"'"$loginJP"'","device":"'"$deviceJP2"'","pk":"'"$pkJP2"'"}' -H "Content-Type: application/json" -H "authorization: Bearer $1" -i
echo -e "\n=========== End Test 2 ===========\nResponse should be HTTP error 201 Created and the summary of the add operation in the body (as JSON)\n\n"

echo -e "=========== Begin Test 3 ===========\nTry to add again a PK ($pkJP1) for the same login and device than before ($loginJP-$deviceJP1)\n"
curl "$baseUrl/public-key" -d '{"login":"'"$loginJP"'","device":"'"$deviceJP1"'","pk":"'"$pkJP1"'"}' -H "Content-Type: application/json" -H "authorization: Bearer $1" -i
echo -e "\n=========== End Test 3 ===========\nResponse should be HTTP error 400 Bad Request\n\n"

echo -e "=========== Begin Test 4 ===========\nTry to get the PK of $loginJP-$deviceJP1\n"
curl "$baseUrl/public-key/$loginJP/$deviceJP1"  -H "authorization: Bearer $1" -i
echo -e "\n=========== End Test 4 ===========\nResponse should be HTTP error 200 OK with PK in response body (as JSON)\n\n"

echo -e "=========== Begin Test 5 ===========\nTry to get the PK of $loginWrong-$deviceWrong\n"
curl "$baseUrl/public-key/$loginWrong/$deviceWrong"  -H "authorization: Bearer $1" -i
echo -e "\n=========== End Test 5 ===========\nResponse should be HTTP error 404 Not Found\n\n"

echo -e "=========== Begin Test 6 ===========\nTry to get all the  PK of $loginJP\n"
curl "$baseUrl/public-key/$loginJP"  -H "authorization: Bearer $1" -i
echo -e "\n=========== End Test 6 ===========\nResponse should be HTTP error 200 OK with all the PK in response body (as JSON)\n\n"

echo -e "=========== Begin Test 7 ===========\nTry to get all the PK of $loginWrong\n"
curl "$baseUrl/public-key/$loginWrong"  -H "authorization: Bearer $1" -i
echo -e "\n=========== End Test 7 ===========\nResponse should be HTTP error 401 Unauthorized\n\n"

pkJP3="pk3"
echo -e "=========== Begin Test 8 ===========\nTry to update the PK of $loginJP\n"
curl "$baseUrl/public-key/$loginJP/$deviceJP1" -X PUT -d '{"pk": "'"$pkJP3"'"}' -H "Content-Type: application/json" -H "authorization: Bearer $1" -i
curl "$baseUrl/public-key/$loginJP"  -H "authorization: Bearer $1" -i
echo -e "\n=========== End Test 8 ===========\nResponse should be HTTP error 200 OK followed by HTTP error 200 OK with PK in response body (as JSON) \n\n"