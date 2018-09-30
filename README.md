# OctoBoard [![Build Status](https://travis-ci.org/radlinskii/octo-board.svg?branch=master)](https://travis-ci.org/radlinskii/octo-board)

## What Is It

**OctoBoard** is a web app created to easily find github issues by labels, organization and language.

App's running on [octo-board.herokuapp.com](https://octo-board.herokuapp.com/).

## How It Works

OctoBoard is consuming [GitHub API v3](https://developer.github.com/v3/) to find issues the user is interested in.
OctoBoard is using [go-github](https://github.com/google/go-github), google's client library to access GitHub API.

User can provide filters to specify what issues he is looking for. There are three text filters:

1. `labels` - comma separated labels applied to issues
2. `language` - main programming language of the repository that contains the filtered issue
3. `organization` - owner of the repository

There are two additional filters:

1. `uncommented` - filters out all the issues that have already got commented
2. `unassigned` - results in displaying only issues that don't have assignees

The results are only **open** issues, sorted by *creation date* in *descending* order.

## Development

These instructions will get you a copy of the project up and running on your local machine for development purposes.

1. clone the repo
  - `git clone https://github.com/radlinskii/octo-board.git`
2. export the `$PORT` enviromental variable
  - `export PORT=3030`
3. run the app inside the repository folder
  - `go run octoboard.go`
  
### Prerequisites

[Go 1.11](https://golang.org/dl/)

## Contributing

Want to contribute? Awesome!:tada:

To fix a bug or add a new feature follow these steps:

1. fork the repository
2. create a new branch
3. Make changes
4. commit your work
5. push your branch to forked repository
6. create a pull request

## Creating Issues

If you have found a bug or want to see a new feature and don't see related issue or pull request just create an issue [here](https://github.com/radlinskii/octo-board/issues/new).
