#!/bin/bash

user=logservice2
group=logservice2

if [ -e '/var/log/GameOperate/client.log' ];then
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

echo '$FileOwner '$user > /etc/rsyslog.d/logservice2.conf
echo '$FileGroup '$user >> /etc/rsyslog.d/logservice2.conf
echo '$FileCreateMode 0644' >> /etc/rsyslog.d/logservice2.conf
echo '$template gamelog,"/var/log/GameOperate/client.log"' >> /etc/rsyslog.d/logservice2.conf
echo 'if ($syslogfacility-text == "local0" or $syslogfacility-text == "LOG_LOCAL0") and $syslogtag contains "cocos" then -?gamelog' >> /etc/rsyslog.d/logservice2.conf
echo '& ~' >> /etc/rsyslog.d/logservice2.conf


service rsyslog restart

if [ -d '/etc/logrotate.d' ];then
	cat>/etc/logrotate.d/logservice2<<EOF
/var/log/GameOperate/client.log{
	daily
	rotate 15
	nocompress
	copytruncate
	dateext
	sharedscripts
	postrotate
		cd /opt/logservice2;./loggie >> ./loggie.log 2>&1 &
	endscript
}
EOF
else
	echo "设置日志回滚失败，没有/etc/logrotate.d目录"
fi
