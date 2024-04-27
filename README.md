# Money Moon API

This project is an API for MoneyMoon, which has all the necessary endpoints for the web to work.

## Models

- Banks
  - \_id
  - name
  - userId
- Groups
  - \_id
  - name
  - userId
  - total
  - transactions `string[]`
- Transactions
  - \_id
  - date
  - dateObject
  - bankId
  - type
  - description
  - amount
  - balance
- Users
  - \_id
  - username
  - email
  - password
  - bank `string []`
  - debts
    - name
    - description
    - amount
    - payed

## Endpoints

| Name                         | Description                                                               | URL                                                                         | Method | Body                                                   | Headers      |
| ---------------------------- | ------------------------------------------------------------------------- | --------------------------------------------------------------------------- | ------ | ------------------------------------------------------ | ------------ |
| Root                         | Checks if the server is running                                           | /                                                                           | GET    | NA                                                     | NA           |
| Transactions by bankId       | Returns all transaction with the given bank id                            | /transactions/{bankId}                                                      | GET    | NA                                                     | access-token |
| Get user by email            | Returns a user by the given email                                         | /user/{email}                                                               | GET    | NA                                                     | access-token |
| Get unpaid debts             | Returns all the unpaid users from a user                                  | /debts                                                                      | GET    | NA                                                     | access-token |
| Get user banks               | Return a list of the user's banks                                         | /user/banks                                                                 | GET    | NA                                                     | access-token |
| Get user groups              | Returns a list of all the groups that the user has                        | /groups                                                                     | GET    | NA                                                     | access-token |
| Search Transactions          | Returns all the transactions that match the given query                   | /transactions/find?search={query}&page={page}&limit={limit}&bankId={bankId} | GET    | NA                                                     | access-token |
| Get transactions by group id | Returns all the transactions as objects from the id list in groups        | /groups/{groupId}                                                           | GET    | NA                                                     | access-token |
| Sign Up                      | Creates a new user account                                                | /signup                                                                     | POST   | `{username, email, password}`                          | NA           |
| Login                        | Checks if the user can login and returns a jwt to set into `access-token` | /login                                                                      | POST   | `{email, password}`                                    | NA           |
| Create bank                  | Creates a new bank for a user                                             | /bank/create                                                                | POST   | `{name, userId}`                                       | access-token |
| Create transactions          | Creates one or more transactions for a user                               | /transactions/create                                                        | POST   | `[{date, bankId, type, description, amount, balance}]` | access-token |
| Create debt                  | Creates a debts for a user                                                | /debts/create                                                               | POST   | `{name, description, Amount}`                          | access-token |
| Create group                 | Creates a new group for the user                                          | /groups                                                                     | POST   | `{name}`                                               | access-token |
| Add transactions to group    | Pushes transactions to a group                                            | /groups/add/{groupId}                                                       | PUT    | `[string]`                                             | access-token |
| Pay debt                     | Substracts the payed amount from the current amount                       | /debts/pay?name={name}&amount={amount}                                      | PUT    | NA                                                     | access-token |
| Delete debt                  | Removes a user debt                                                       | /debts/remove/{name}                                                        | DELETE | NA                                                     | access-token |
| Delete group transaction     | Removes a transaction form a group                                        | /groups/delete/{groupId}                                                    | DELETE | `{transactionId}`                                      | access-token |
