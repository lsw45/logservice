package errorx

type ErrInt uint32

const ErrNil ErrInt = 0

//全局错误码
const SERVER_COMMON_ERROR ErrInt = 100001
const REUQEST_PARAM_ERROR ErrInt = 100002
const DB_ERROR ErrInt = 100003
const AUTH_ERROR ErrInt = 100004
const DeadlineExceeded_ERROR ErrInt = 100005
const HTTP_REQUEST_ERR ErrInt = 100006
const RPC_TIME_OUT ErrInt = 100007
const NOT_FOUNT_RETURN_NULL ErrInt = 100008

//监控模块错误码
const INSTANCE_PARAMS_ERR ErrInt = 200001
const METRIC_NOTFOUND_ERR ErrInt = 200002
const LEVEL_NOTFOUND_ERR ErrInt = 200003
const GROUP_NOTFOUND_ERR ErrInt = 200004
const OP_CROSS_BORDER ErrInt = 200005
const OUT_RANGE_METRIC ErrInt = 200006
const APPLY_CONFIG_ERR ErrInt = 200007
const PROME_NOT_FOUND ErrInt = 200008
const PROME_OUT_RANGE ErrInt = 200009
const EXPORTER_NOT_FOUND ErrInt = 2000010
const INSTANCE_ID_CAN_NOT_NIL ErrInt = 2000011
const PARAM_ERR_PID_NIL ErrInt = 2000012
const CUSTOM_PROCESS_MONITOR_TOO_MANY ErrInt = 2000013
const CREATE_TASK_ERR ErrInt = 2000014
const GET_TOP_TASK_ERR ErrInt = 2000015
const INIT_HTTP_ERR ErrInt = 2000016
const SALT_API_RESP_DECODE_ERR ErrInt = 2000017
const TOP_PROCESS_NOT_FOUND ErrInt = 2000018
const TOP_PROCESS_MONITOR_TOO_MANY ErrInt = 2000019
const TOP_PROCESS_MONITOR_REPEATE ErrInt = 2000020
const SALT_EXPORTER_ERR ErrInt = 2000021
const PROCESS_NUM_RULE_ERR ErrInt = 2000022
const UPDATE_RULE_TASK_ERR ErrInt = 2000023
const UPDATE_RECEIVER_TASK_ERR ErrInt = 2000024
const UPDATE_PROCESS_TASK_ERR ErrInt = 2000025
const RULE_NOT_FOUND ErrInt = 2000026
const PROCESS_NAME_ERR ErrInt = 2000027
const NODE_LIST_STEP_TO_LARGE ErrInt = 2000028
const MESSAGE_CONSUME_ERR ErrInt = 2000029
const REDIS_LOCK_FAIL ErrInt = 2000030
const MESSAGE_REPEATE ErrInt = 2000031
const DB_INSTANCE_NOT_FOUND ErrInt = 2000032
const MONGO_QUERY_DATE_ERR ErrInt = 2000033

var message map[ErrInt]string

func init() {
	message = make(map[ErrInt]string)
	message[SERVER_COMMON_ERROR] = "服务器开小差啦,稍后再来试一试"
	message[REUQEST_PARAM_ERROR] = "参数错误"
	message[DB_ERROR] = "数据库繁忙,请稍后再试"
	message[AUTH_ERROR] = "鉴权失败"
	message[DeadlineExceeded_ERROR] = "请求超时"
	message[HTTP_REQUEST_ERR] = "调用外部接口失败"
	message[RPC_TIME_OUT] = "调用服务超时"
	message[NOT_FOUNT_RETURN_NULL] = "数据不存在"

	//监控模块
	message[INSTANCE_PARAMS_ERR] = "实例地址参数格式错误"
	message[METRIC_NOTFOUND_ERR] = "选择指标不存在"
	message[LEVEL_NOTFOUND_ERR] = "选择监控等级不存在"
	message[GROUP_NOTFOUND_ERR] = "联系人组数据错误"
	message[OP_CROSS_BORDER] = "操作越权"
	message[OUT_RANGE_METRIC] = "监控指标数量超过限制"
	message[APPLY_CONFIG_ERR] = "应用配置失败"
	message[PROME_NOT_FOUND] = "未找到对应监控服务"
	message[PROME_OUT_RANGE] = "查询范围过长"
	message[EXPORTER_NOT_FOUND] = "未找到监控实例"
	message[INSTANCE_ID_CAN_NOT_NIL] = "监控实例id不可为空"
	message[PARAM_ERR_PID_NIL] = "PID不可为空"
	message[CUSTOM_PROCESS_MONITOR_TOO_MANY] = "自定义监控进程数量不能超过20"
	message[CREATE_TASK_ERR] = "创建任务失败"
	message[GET_TOP_TASK_ERR] = "获取top任务失败"
	message[INIT_HTTP_ERR] = "初始化http失败"
	message[SALT_API_RESP_DECODE_ERR] = "saltApi数据解析失败"
	message[TOP_PROCESS_NOT_FOUND] = "未找到对应进程"
	message[TOP_PROCESS_MONITOR_TOO_MANY] = "最多开启20个进程监控"
	message[TOP_PROCESS_MONITOR_REPEATE] = "不可重复创建"
	message[SALT_EXPORTER_ERR] = "执行监控脚本失败"
	message[PROCESS_NUM_RULE_ERR] = "进程数量监控参数缺失"
	message[UPDATE_RULE_TASK_ERR] = "更新告警规则任务失败"
	message[UPDATE_RECEIVER_TASK_ERR] = "更新联系人任务失败"
	message[UPDATE_PROCESS_TASK_ERR] = "更新进程监控任务失败"
	message[RULE_NOT_FOUND] = "未找到对于报警规则"
	message[PROCESS_NAME_ERR] = "进程名称包含非法字符"
	message[NODE_LIST_STEP_TO_LARGE] = "时间间隔过大"
	message[MESSAGE_CONSUME_ERR] = "消息无法消费"
	message[REDIS_LOCK_FAIL] = "操作冲突，请重试"
	message[MESSAGE_REPEATE] = "消息重复"
	message[DB_INSTANCE_NOT_FOUND] = "未找到数据库实例"
	message[MONGO_QUERY_DATE_ERR] = "时间格式错误"

}

func MapErrMsg(errcode ErrInt) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

func IsCodeErr(errcode ErrInt) bool {
	if _, ok := message[errcode]; ok {
		return true
	} else {
		return false
	}
}
