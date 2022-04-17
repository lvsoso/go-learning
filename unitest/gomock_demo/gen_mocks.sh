#!/usr/bin/env bash

mockgen -source=db.go -destination=mocks/db_mock.go -package=mocks
