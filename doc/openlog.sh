#!/bin/bash

user=logservice2
group=logservice2

if [ -e '/var/log/gameOperate/client.log' ];then
	echo "运维日志重定向已启动，可直接采集"
	exit 0
fi

groupadd -r -g 1001 -o $group
if [ $? -ne 0 ];then
	echo "新建用户组失败"
	exit 1
fi

useradd -r -g $group $user

if [ $? -ne 0 ];then
	echo "新建用户失败"
	exit 1
fi

echo '$FileOwner '$user > /etc/rsyslog.d/gameOperator.conf
echo '$FileGroup '$user >> /etc/rsyslog.d/gameOperator.conf
echo '$FileCreateMode 0644' >> /etc/rsyslog.d/gameOperator.conf
echo 'template(name="ctpl" type="string" string="%TIMESTAMP:::date-unixtimestamp% %programname% %msg%\n")' >> /etc/rsyslog.d/gameOperator.conf
echo '$template gamelog,"/var/log/gameOperate/client.log"' >> /etc/rsyslog.d/gameOperator.conf
echo 'if ($syslogfacility-text == "local0" or $syslogfacility-text == "LOG_LOCAL0") and $syslogtag contains "cocos" then -?gamelog;ctpl' >> /etc/rsyslog.d/gameOperator.conf
echo '& ~' >> /etc/rsyslog.d/gameOperator.conf

echo '$FileOwner '$user > /etc/rsyslog.d/gameServer.conf
echo '$FileGroup '$user >> /etc/rsyslog.d/gameServer.conf
echo '$FileCreateMode 0644' >> /etc/rsyslog.d/gameServer.conf
echo 'template(name="stpl" type="string" string="%TIMESTAMP:::date-unixtimestamp% %programname% %msg%\n")' >> /etc/rsyslog.d/gameServer.conf
echo '$template serverlog,"/var/log/engine/server.log"' >> /etc/rsyslog.d/gameServer.conf
echo 'if $syslogtag contains "supervisord" then -?serverlog;stpl' >> /etc/rsyslog.d/gameServer.conf
echo '& ~' >> /etc/rsyslog.d/gameServer.conf

service rsyslog restart

if [ -d '/etc/logrotate.d' ];then
	cat>/etc/logrotate.d/gameOperator<<EOF
/var/log/gameOperate/client.log{
	daily
	rotate 15
	nocompress
	copytruncate
	dateext
	sharedscripts
	postrotate
		systemctl restart loggie
	endscript
}
EOF
else
	echo "设置日志回滚失败，没有/etc/logrotate.d目录"
fi
