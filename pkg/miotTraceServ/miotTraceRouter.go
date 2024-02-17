package miottraceserv

import mttypes "miot_tracing_go/mtTypes"

// router
func (m *MiotTracingServImpl) handleMessage(message mttypes.Message) error {
	switch message.Type {
	case mttypes.TYPE_FIRST_UPLOAD:
		err := m.handleFirstData(message)
		if err != nil {
			iotlog.Errorf("handleFirstData error: %v", err)
		}
	case mttypes.TYPE_UPLOAD_THIRD_INDEX:
		err := m.handleUploadThirdIndex(message)
		if err != nil {
			iotlog.Errorf("handleUploadThirdIndex error: %v", err)
		}
	case mttypes.TYPE_UPLOAD_TAXI_DATA:
		err := m.handleData(message)
		if err != nil {
			iotlog.Errorf("handleData error: %v", err)
		}
	case mttypes.TYPE_UPDATE_SECOND_INDEX:
		err := m.handleUpdateSecondIndex(message)
		if err != nil {
			iotlog.Errorf("handleUpdateSecondIndex error: %v", err)
		}
	case mttypes.TYPE_UPDATE_THIRD_INDEX:
		err := m.handleUpdateThirdIndex(message)
		if err != nil {
			iotlog.Errorf("handleUpdateThirdIndex error: %v", err)
		}
	case mttypes.TYPE_UPLOAD_THIRD_INDEX_HEAD:
		err := m.handleUploadMetaData(message)
		if err != nil {
			iotlog.Errorf("handleUploadMetaData error: %v", err)
		}
	case mttypes.TYPE_BUILD_QUERY:
		err := m.handleBuildQuery(message)
		if err != nil {
			iotlog.Errorf("handleQuery error: %v", err)
		}
	case mttypes.TYPE_SEARCH_THIRD_INDEX:
		err := m.handleSearchThirdIndex(message)
		if err != nil {
			iotlog.Errorf("handleSearchThirdIndex error: %v", err)
		}
	case mttypes.TYPE_QUERY_TAXI_DATA:
		err := m.handleQueryData(message)
		if err != nil {
			iotlog.Errorf("handleQueryData error: %v", err)
		}
	default:
		iotlog.Errorf("unknown message type: %v", message.Type)
	}
	return nil
}
