#!/bin/bash

.PHONY: run-http
run-http:
	@go run ./sources/http

.PHONY: run-workers
run-workers:
	@go run ./sources/worker
