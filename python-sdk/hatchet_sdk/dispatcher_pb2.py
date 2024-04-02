# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: dispatcher.proto
# Protobuf Python Version: 4.25.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x10\x64ispatcher.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"p\n\x15WorkerRegisterRequest\x12\x12\n\nworkerName\x18\x01 \x01(\t\x12\x0f\n\x07\x61\x63tions\x18\x02 \x03(\t\x12\x10\n\x08services\x18\x03 \x03(\t\x12\x14\n\x07maxRuns\x18\x04 \x01(\x05H\x00\x88\x01\x01\x42\n\n\x08_maxRuns\"P\n\x16WorkerRegisterResponse\x12\x10\n\x08tenantId\x18\x01 \x01(\t\x12\x10\n\x08workerId\x18\x02 \x01(\t\x12\x12\n\nworkerName\x18\x03 \x01(\t\"\x84\x02\n\x0e\x41ssignedAction\x12\x10\n\x08tenantId\x18\x01 \x01(\t\x12\x15\n\rworkflowRunId\x18\x02 \x01(\t\x12\x18\n\x10getGroupKeyRunId\x18\x03 \x01(\t\x12\r\n\x05jobId\x18\x04 \x01(\t\x12\x0f\n\x07jobName\x18\x05 \x01(\t\x12\x10\n\x08jobRunId\x18\x06 \x01(\t\x12\x0e\n\x06stepId\x18\x07 \x01(\t\x12\x11\n\tstepRunId\x18\x08 \x01(\t\x12\x10\n\x08\x61\x63tionId\x18\t \x01(\t\x12\x1f\n\nactionType\x18\n \x01(\x0e\x32\x0b.ActionType\x12\x15\n\ractionPayload\x18\x0b \x01(\t\x12\x10\n\x08stepName\x18\x0c \x01(\t\"\'\n\x13WorkerListenRequest\x12\x10\n\x08workerId\x18\x01 \x01(\t\",\n\x18WorkerUnsubscribeRequest\x12\x10\n\x08workerId\x18\x01 \x01(\t\"?\n\x19WorkerUnsubscribeResponse\x12\x10\n\x08tenantId\x18\x01 \x01(\t\x12\x10\n\x08workerId\x18\x02 \x01(\t\"\xe1\x01\n\x13GroupKeyActionEvent\x12\x10\n\x08workerId\x18\x01 \x01(\t\x12\x15\n\rworkflowRunId\x18\x02 \x01(\t\x12\x18\n\x10getGroupKeyRunId\x18\x03 \x01(\t\x12\x10\n\x08\x61\x63tionId\x18\x04 \x01(\t\x12\x32\n\x0e\x65ventTimestamp\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12+\n\teventType\x18\x06 \x01(\x0e\x32\x18.GroupKeyActionEventType\x12\x14\n\x0c\x65ventPayload\x18\x07 \x01(\t\"\xec\x01\n\x0fStepActionEvent\x12\x10\n\x08workerId\x18\x01 \x01(\t\x12\r\n\x05jobId\x18\x02 \x01(\t\x12\x10\n\x08jobRunId\x18\x03 \x01(\t\x12\x0e\n\x06stepId\x18\x04 \x01(\t\x12\x11\n\tstepRunId\x18\x05 \x01(\t\x12\x10\n\x08\x61\x63tionId\x18\x06 \x01(\t\x12\x32\n\x0e\x65ventTimestamp\x18\x07 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\'\n\teventType\x18\x08 \x01(\x0e\x32\x14.StepActionEventType\x12\x14\n\x0c\x65ventPayload\x18\t \x01(\t\"9\n\x13\x41\x63tionEventResponse\x12\x10\n\x08tenantId\x18\x01 \x01(\t\x12\x10\n\x08workerId\x18\x02 \x01(\t\"9\n SubscribeToWorkflowEventsRequest\x12\x15\n\rworkflowRunId\x18\x01 \x01(\t\"\xb2\x02\n\rWorkflowEvent\x12\x15\n\rworkflowRunId\x18\x01 \x01(\t\x12#\n\x0cresourceType\x18\x02 \x01(\x0e\x32\r.ResourceType\x12%\n\teventType\x18\x03 \x01(\x0e\x32\x12.ResourceEventType\x12\x12\n\nresourceId\x18\x04 \x01(\t\x12\x32\n\x0e\x65ventTimestamp\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x14\n\x0c\x65ventPayload\x18\x06 \x01(\t\x12\x0e\n\x06hangup\x18\x07 \x01(\x08\x12\x18\n\x0bstepRetries\x18\x08 \x01(\x05H\x00\x88\x01\x01\x12\x17\n\nretryCount\x18\t \x01(\x05H\x01\x88\x01\x01\x42\x0e\n\x0c_stepRetriesB\r\n\x0b_retryCount\"W\n\rOverridesData\x12\x11\n\tstepRunId\x18\x01 \x01(\t\x12\x0c\n\x04path\x18\x02 \x01(\t\x12\r\n\x05value\x18\x03 \x01(\t\x12\x16\n\x0e\x63\x61llerFilename\x18\x04 \x01(\t\"\x17\n\x15OverridesDataResponse\"U\n\x10HeartbeatRequest\x12\x10\n\x08workerId\x18\x01 \x01(\t\x12/\n\x0bheartbeatAt\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\"\x13\n\x11HeartbeatResponse*N\n\nActionType\x12\x12\n\x0eSTART_STEP_RUN\x10\x00\x12\x13\n\x0f\x43\x41NCEL_STEP_RUN\x10\x01\x12\x17\n\x13START_GET_GROUP_KEY\x10\x02*\xa2\x01\n\x17GroupKeyActionEventType\x12 \n\x1cGROUP_KEY_EVENT_TYPE_UNKNOWN\x10\x00\x12 \n\x1cGROUP_KEY_EVENT_TYPE_STARTED\x10\x01\x12\"\n\x1eGROUP_KEY_EVENT_TYPE_COMPLETED\x10\x02\x12\x1f\n\x1bGROUP_KEY_EVENT_TYPE_FAILED\x10\x03*\x8a\x01\n\x13StepActionEventType\x12\x1b\n\x17STEP_EVENT_TYPE_UNKNOWN\x10\x00\x12\x1b\n\x17STEP_EVENT_TYPE_STARTED\x10\x01\x12\x1d\n\x19STEP_EVENT_TYPE_COMPLETED\x10\x02\x12\x1a\n\x16STEP_EVENT_TYPE_FAILED\x10\x03*e\n\x0cResourceType\x12\x19\n\x15RESOURCE_TYPE_UNKNOWN\x10\x00\x12\x1a\n\x16RESOURCE_TYPE_STEP_RUN\x10\x01\x12\x1e\n\x1aRESOURCE_TYPE_WORKFLOW_RUN\x10\x02*\xfe\x01\n\x11ResourceEventType\x12\x1f\n\x1bRESOURCE_EVENT_TYPE_UNKNOWN\x10\x00\x12\x1f\n\x1bRESOURCE_EVENT_TYPE_STARTED\x10\x01\x12!\n\x1dRESOURCE_EVENT_TYPE_COMPLETED\x10\x02\x12\x1e\n\x1aRESOURCE_EVENT_TYPE_FAILED\x10\x03\x12!\n\x1dRESOURCE_EVENT_TYPE_CANCELLED\x10\x04\x12!\n\x1dRESOURCE_EVENT_TYPE_TIMED_OUT\x10\x05\x12\x1e\n\x1aRESOURCE_EVENT_TYPE_STREAM\x10\x06\x32\xd1\x04\n\nDispatcher\x12=\n\x08Register\x12\x16.WorkerRegisterRequest\x1a\x17.WorkerRegisterResponse\"\x00\x12\x33\n\x06Listen\x12\x14.WorkerListenRequest\x1a\x0f.AssignedAction\"\x00\x30\x01\x12\x35\n\x08ListenV2\x12\x14.WorkerListenRequest\x1a\x0f.AssignedAction\"\x00\x30\x01\x12\x34\n\tHeartbeat\x12\x11.HeartbeatRequest\x1a\x12.HeartbeatResponse\"\x00\x12R\n\x19SubscribeToWorkflowEvents\x12!.SubscribeToWorkflowEventsRequest\x1a\x0e.WorkflowEvent\"\x00\x30\x01\x12?\n\x13SendStepActionEvent\x12\x10.StepActionEvent\x1a\x14.ActionEventResponse\"\x00\x12G\n\x17SendGroupKeyActionEvent\x12\x14.GroupKeyActionEvent\x1a\x14.ActionEventResponse\"\x00\x12<\n\x10PutOverridesData\x12\x0e.OverridesData\x1a\x16.OverridesDataResponse\"\x00\x12\x46\n\x0bUnsubscribe\x12\x19.WorkerUnsubscribeRequest\x1a\x1a.WorkerUnsubscribeResponse\"\x00\x42GZEgithub.com/hatchet-dev/hatchet/internal/services/dispatcher/contractsb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'dispatcher_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  _globals['DESCRIPTOR']._options = None
  _globals['DESCRIPTOR']._serialized_options = b'ZEgithub.com/hatchet-dev/hatchet/internal/services/dispatcher/contracts'
  _globals['_ACTIONTYPE']._serialized_start=1780
  _globals['_ACTIONTYPE']._serialized_end=1858
  _globals['_GROUPKEYACTIONEVENTTYPE']._serialized_start=1861
  _globals['_GROUPKEYACTIONEVENTTYPE']._serialized_end=2023
  _globals['_STEPACTIONEVENTTYPE']._serialized_start=2026
  _globals['_STEPACTIONEVENTTYPE']._serialized_end=2164
  _globals['_RESOURCETYPE']._serialized_start=2166
  _globals['_RESOURCETYPE']._serialized_end=2267
  _globals['_RESOURCEEVENTTYPE']._serialized_start=2270
  _globals['_RESOURCEEVENTTYPE']._serialized_end=2524
  _globals['_WORKERREGISTERREQUEST']._serialized_start=53
  _globals['_WORKERREGISTERREQUEST']._serialized_end=165
  _globals['_WORKERREGISTERRESPONSE']._serialized_start=167
  _globals['_WORKERREGISTERRESPONSE']._serialized_end=247
  _globals['_ASSIGNEDACTION']._serialized_start=250
  _globals['_ASSIGNEDACTION']._serialized_end=510
  _globals['_WORKERLISTENREQUEST']._serialized_start=512
  _globals['_WORKERLISTENREQUEST']._serialized_end=551
  _globals['_WORKERUNSUBSCRIBEREQUEST']._serialized_start=553
  _globals['_WORKERUNSUBSCRIBEREQUEST']._serialized_end=597
  _globals['_WORKERUNSUBSCRIBERESPONSE']._serialized_start=599
  _globals['_WORKERUNSUBSCRIBERESPONSE']._serialized_end=662
  _globals['_GROUPKEYACTIONEVENT']._serialized_start=665
  _globals['_GROUPKEYACTIONEVENT']._serialized_end=890
  _globals['_STEPACTIONEVENT']._serialized_start=893
  _globals['_STEPACTIONEVENT']._serialized_end=1129
  _globals['_ACTIONEVENTRESPONSE']._serialized_start=1131
  _globals['_ACTIONEVENTRESPONSE']._serialized_end=1188
  _globals['_SUBSCRIBETOWORKFLOWEVENTSREQUEST']._serialized_start=1190
  _globals['_SUBSCRIBETOWORKFLOWEVENTSREQUEST']._serialized_end=1247
  _globals['_WORKFLOWEVENT']._serialized_start=1250
  _globals['_WORKFLOWEVENT']._serialized_end=1556
  _globals['_OVERRIDESDATA']._serialized_start=1558
  _globals['_OVERRIDESDATA']._serialized_end=1645
  _globals['_OVERRIDESDATARESPONSE']._serialized_start=1647
  _globals['_OVERRIDESDATARESPONSE']._serialized_end=1670
  _globals['_HEARTBEATREQUEST']._serialized_start=1672
  _globals['_HEARTBEATREQUEST']._serialized_end=1757
  _globals['_HEARTBEATRESPONSE']._serialized_start=1759
  _globals['_HEARTBEATRESPONSE']._serialized_end=1778
  _globals['_DISPATCHER']._serialized_start=2527
  _globals['_DISPATCHER']._serialized_end=3120
# @@protoc_insertion_point(module_scope)
