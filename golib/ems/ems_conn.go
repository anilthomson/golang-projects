package ems

/*
#cgo CFLAGS: -I /opt/tibco/ems/8.4/include/tibems
#cgo LDFLAGS: -L/opt/tibco/ems/8.4/lib -L/opt/tibco/ems/8.4/lib/64  -ltibems64
#include <tibems.h>
#include <confact.h>
#include <emsadmin.h>
#include <dest.h>
#include "ems.h"




*/
import "C"
import (
	"fmt"
	"unsafe"
)

func checkStatus(status *C.tibems_status, C_errorContext *C.tibemsErrorContext) {
	if *status != 0 {
		var errCode = C.GoString(C.tibemsStatus_GetText(*status))
		if C_errorContext != nil {
			var err *C.char
			C.tibemsErrorContext_GetLastErrorString(*C_errorContext, &err)
			var ErrorText = C.GoString(err)
			C.free(unsafe.Pointer(err))
			fmt.Println(errCode)
			if len(ErrorText) != 0 {
				charErr := errCode + " - " + ErrorText[20:]
				panic(charErr)
			} else {
				panic(errCode)
			}
		}
	}
}

func CreateConnection(emsServerURL string, username string, password string) Connection {
	var connection jmsConnection
	var C_errorContext C.tibemsErrorContext
	status := C.tibemsErrorContext_Create(&C_errorContext)
	var c_factory C.tibemsConnectionFactory = C.tibemsConnectionFactory_Create()
	serverUrl := C.CString(emsServerURL)
	status = C.tibemsConnectionFactory_SetServerURL(c_factory, serverUrl)
	checkStatus(&status, &C_errorContext)
	status = C.tibemsConnectionFactory_CreateConnection(c_factory, &connection.C_connection, C.CString(username), C.CString(password))
	checkStatus(&status, &C_errorContext)
	C.tibemsErrorContext_Close(C_errorContext)
	return &connection
}

func (connection *jmsConnection) CreateSession(transacted bool, ack int) Session {
	var session jmsSession
	/* Each thread has it's own error context. */
	var C_errorContext C.tibemsErrorContext
	status := C.tibemsErrorContext_Create(&C_errorContext)
	checkStatus(&status, nil)
	status = C.tibemsConnection_CreateSession(connection.C_connection, &session.C_session, C.TIBEMS_FALSE, C.TIBEMS_AUTO_ACKNOWLEDGE)
	checkStatus(&status, session.ErrorContext)
	session.ErrorContext = &C_errorContext
	return &session
}
func (connection *jmsConnection) Start() {
	status := C.tibemsConnection_Start(connection.C_connection)
	checkStatus(&status, nil)
}
func (connection *jmsConnection) Stop() {
	status := C.tibemsConnection_Stop(connection.C_connection)
	checkStatus(&status, nil)
}
func (connection *jmsConnection) Close() {
	status := C.tibemsConnection_Close(connection.C_connection)
	checkStatus(&status, nil)
}
