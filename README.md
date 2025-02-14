# Go - Google Sheets Reader

This server allows for reading of a Google Sheet.

## Setup

1. On Google Console, enable the Google Sheets API and setup a service account.
2. Generate the JSON credentials for the service account and put these in the `.env` as seen in `.env.example`
3. Go to the Google Sheet that you want to use, and share it with the email of the Service Account.
4. Extract the ID of the sheet from the URL.

## Development

- Copy `.env.example` to `.env` and fill in the details.
- To Run;
  - If using `air`, can run `air ./api/cmd` for hot reloading
  - Else can use `go run ./api/cmd` and then rerun every time a change occurs
