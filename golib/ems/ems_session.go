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
import "unsafe"

func (session *jmsSession) CreateQueueProducer(queuename string) MessageProducer {
	var msgProducer jmsMessageProducer
	var destination C.tibemsDestination
	var queue C.tibemsQueue
	status := C.tibemsQueue_Create(&queue, C.CString(queuename))
	checkStatus(&status, session.ErrorContext)
	destination = _Ctype_tibemsDestination(queue)
	status = C.tibemsSession_CreateProducer(session.C_session, &msgProducer.C_msgProducer, destination)
	checkStatus(&status, session.ErrorContext)
	msgProducer.C_destination = destination
	msgProducer.ErrorContext = session.ErrorContext
	return &msgProducer
}
func (session *jmsSession) CreateTopicProducer(topicname string) MessageProducer {
	var msgProducer jmsMessageProducer
	var destination C.tibemsDestination
	var topic C.tibemsTopic
	Ctopicname := C.CString(topicname)
	status := C.tibemsTopic_Create(&topic, Ctopicname)
	checkStatus(&status, session.ErrorContext)
	destination = _Ctype_tibemsDestination(topic)
	status = C.tibemsSession_CreateProducer(session.C_session, &msgProducer.C_msgProducer, destination)
	checkStatus(&status, session.ErrorContext)
	msgProducer.ErrorContext = session.ErrorContext
	C.free(unsafe.Pointer(Ctopicname))
	return &msgProducer
}
func (session *jmsSession) CreateQueueConsumer(queuename string) MessageConsumer {
	var msgConsumer jmsMessageConsumer
	var destination C.tibemsDestination
	var queue C.tibemsQueue
	status := C.tibemsQueue_Create(&queue, C.CString(queuename))
	checkStatus(&status, session.ErrorContext)
	destination = _Ctype_tibemsDestination(queue)
	status = C.tibemsSession_CreateConsumer(session.C_session, &msgConsumer.C_msgConsumer, destination, nil, C.TIBEMS_FALSE)
	checkStatus(&status, session.ErrorContext)
	msgConsumer.ErrorContext = session.ErrorContext
	return &msgConsumer
}

func (session *jmsSession) CreateTopicConsumer(topicname string) MessageConsumer {
	var msgConsumer jmsMessageConsumer
	var destination C.tibemsDestination
	var topic C.tibemsTopic
	status := C.tibemsTopic_Create(&topic, C.CString(topicname))
	checkStatus(&status, session.ErrorContext)
	destination = _Ctype_tibemsDestination(topic)
	status = C.tibemsSession_CreateConsumer(session.C_session, &msgConsumer.C_msgConsumer, destination, nil, C.TIBEMS_FALSE)
	checkStatus(&status, session.ErrorContext)
	msgConsumer.ErrorContext = session.ErrorContext
	return &msgConsumer
}

func (session *jmsSession) Close() {
	status := C.tibemsSession_Close(session.C_session)
	checkStatus(&status, session.ErrorContext)
	status = C.tibemsErrorContext_Close(*session.ErrorContext)
	checkStatus(&status, nil)
}
