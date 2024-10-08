# KaraBank
## About this app
Welcome to KaraBank, your bank of trust ;)

## How to run the app
**Prerequisites:**
- Docker installed (I am using Docker Desktop version 4.33.1 on Windows 11)

**Installation steps:**
- Clone the code with `https://github.com/karaMuha/kara-bank.git`
- Inside the root directory of the project run the command `make setup` (this will create the folder db-data to persist data from the postgres container)
- run `make start` (the database will be initialized with the script `init.sql`. You can find the script in the folder `db-script`)
- the http server will listen on port 8080 and postgres on port 5433

## Usage
- POST /users/register -> Register as a customer of our trustworthy bank.
```
{
    "email": "test@test.com",
    "password": "test1234",
    "first_name": "Max",
    "last_name": "Mustermann"
}
```
- POST /users/login -> Login with your credentials.
```
{
    "email": "test@test.com",
    "password": "test1234"
}
```
- POST /accounts -> Create a bank account in order to become rich. Need to be logged in to do so.
```
{
    "currency": "EUR"
}
```
- GET /accounts/{id} -> Get account with provided id. Admin and Banker role can get any account. Customer role can only get his own accoutns.
- GET /accounts -> Admin and Banker role can list accounts.
```
{
    "limit": {any number >= 1},
    "offset": {any number >= 0}
}
```
- POST /transfers -> Transfer money from one account to another. Need to be logged in and you can only send money from your own account.
```
{
    "from_account_id": {id of a created account},
    "to_account_id": {id of another created account},
    "amount": {any number}
}
```

## ToDos
- implement money deposit feature
- implement money withdraw feature
- implement individual lower limit for bank account feature
- implement currency conversion for transactions between accounts with different currencies