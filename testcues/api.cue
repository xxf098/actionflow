package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {
		api: core.#API & {
			host: "https://postman-echo.com"
            method: "GET"
            path: "/get"
            query: {
                cow: "moo"
                task: "task1"
            }
            resp: { body: string }
		}

        print: core.#Stdout & {
            text: api.resp.body
        }

	}
}
