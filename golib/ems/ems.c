#include <tibems.h>
#include <confact.h>
#include <emsadmin.h>
#include <dest.h>
#include "ems.h"
void* ptr;
extern void OnEMSMessage(tibemsMsg msg ,tibemsMsgConsumer consumer);
 

void onMessage(tibemsMsgConsumer consumer, tibemsMsg msg )
{
    OnEMSMessage(msg,consumer);
 }