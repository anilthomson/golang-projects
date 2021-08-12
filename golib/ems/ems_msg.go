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
	"unsafe"
)
import "C"

var listenerMap map[C.tibemsMsgConsumer]*jmsMessageConsumer = make(map[C.tibemsMsgConsumer]*jmsMessageConsumer)

//export OnEMSMessage
func OnEMSMessage(msg C.tibemsMsg, consumer C.tibemsMsgConsumer) {
	var charTxt *_Ctype_char
	status := C.tibemsTextMsg_GetText(msg, &charTxt)
	checkStatus(&status, nil)
	var text = C.GoString(charTxt)
	jmsConsumer := listenerMap[consumer]
	message := &jmsMessage{msg, text, nil, nil}
	jmsConsumer.Listener(message)
}

func (msgProducer *jmsMessageProducer) Send(txtMsg TextMessage) {
	var jmsTxtMsg *jmsTextMessage = txtMsg.(*jmsTextMessage)
	status := C.tibemsMsgProducer_Send(msgProducer.C_msgProducer, jmsTxtMsg.C_txtMsg)
	checkStatus(&status, msgProducer.ErrorContext)
}
func (msgProducer *jmsMessageProducer) Close() {
	status := C.tibemsMsgProducer_Close(msgProducer.C_msgProducer)
	checkStatus(&status, msgProducer.ErrorContext)
	status = C.tibemsDestination_Destroy(msgProducer.C_destination)
	checkStatus(&status, msgProducer.ErrorContext)
}

func (msgConsumer *jmsMessageConsumer) Receive() string {
	var msg C.tibemsMsg
	var charTxt *_Ctype_char
	status := C.tibemsMsgConsumer_Receive(msgConsumer.C_msgConsumer, &msg)
	checkStatus(&status, msgConsumer.ErrorContext)
	status = C.tibemsTextMsg_GetText(msg, &charTxt)
	checkStatus(&status, msgConsumer.ErrorContext)
	var text = C.GoString(charTxt)
	return text
}

func (msgConsumer *jmsMessageConsumer) SetMsgListener(messageListener MessageListener) {

	callback := C.tibemsMsgCallback(unsafe.Pointer(C.onMessage))
	msgConsumer.Listener = messageListener
	msgConsumer.Callback = callback
	status := C.tibemsMsgConsumer_SetMsgListener(msgConsumer.C_msgConsumer, callback, nil)
	checkStatus(&status, msgConsumer.ErrorContext)
	listenerMap[msgConsumer.C_msgConsumer] = msgConsumer

}
func (msg *jmsMessage) GetText() string {

	return msg.text
}
func (msg *jmsTextMessage) Destroy() {
	status := C.tibemsMsg_Destroy(msg.C_txtMsg)
	checkStatus(&status, msg.ErrorContext)
}

// CreateTextMessage ... sss
func CreateTextMessage(message string) TextMessage {
	var txtMsg jmsTextMessage
	status := C.tibemsTextMsg_Create(&txtMsg.C_txtMsg)
	checkStatus(&status, nil)
	var C_text *_Ctype_char = C.CString(message)
	status = C.tibemsTextMsg_SetText(txtMsg.C_txtMsg, C_text)
	checkStatus(&status, nil)
	C.free(unsafe.Pointer(C_text))
	return &txtMsg
}
