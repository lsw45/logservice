package infra

import (
	"context"
	"errors"
	"log-ext/common"
	"log-ext/domain/entity"

	"github.com/olivere/elastic/v7"
)

/*
{
    "total":{
        "value":28,
        "relation":"eq"
    },
    "max_score":1,
    "hits":[
        {
            "_index":"test_index",
            "_id":"gNbqY4MB5bsVxw97I68H",
            "_score":1,
            "_ignored":[
                "message.keyword",
                "event.original.keyword"
            ],
            "_source":{
                "message":"{\"state\":{\"hostname\":\"paas-dev\",\"pipeline\":\"demo\",\"source\":\"kdump\",\"filename\":\"/var/log/kdump.log\",\"timestamp\":\"2022-09-22T14:37:08.689Z\",\"offset\":0,\"bytes\":559},\"fields\":{\"service\":\"/usr/local/go\"},\"body\":\"+ 2022-06-27 16:54:57 /usr/bin/kdumpctl@708: /sbin/kexec -s -d -p '--command-line=BOOT_IMAGE=(hd0,msdos1)/boot/vmlinuz-4.18.0-348.7.1.el8_5.x86_64 ro net.ifnames=0 consoleblank=600 console=tty0 console=ttyS0,115200n8 spectre_v2=off nopti irqpoll nr_cpus=1 reset_devices cgroup_disable=memory mce=off numa=off udev.children-max=2 panic=10 rootflags=nofail acpi_no_memhotplug transparent_hugepage=never nokaslr novmcoredd hest_disable disable_cpu_apicid=0' --initrd=/boot/initramfs-4.18.0-348.7.1.el8_5.x86_64kdump.img /boot/vmlinuz-4.18.0-348.7.1.el8_5.x86_64\"}",
                "@timestamp":"2022-09-22T06:37:11.724606454Z",
                "event":{
                    "original":"{\"state\":{\"hostname\":\"paas-dev\",\"pipeline\":\"demo\",\"source\":\"kdump\",\"filename\":\"/var/log/kdump.log\",\"timestamp\":\"2022-09-22T14:37:08.689Z\",\"offset\":0,\"bytes\":559},\"fields\":{\"service\":\"/usr/local/go\"},\"body\":\"+ 2022-06-27 16:54:57 /usr/bin/kdumpctl@708: /sbin/kexec -s -d -p '--command-line=BOOT_IMAGE=(hd0,msdos1)/boot/vmlinuz-4.18.0-348.7.1.el8_5.x86_64 ro net.ifnames=0 consoleblank=600 console=tty0 console=ttyS0,115200n8 spectre_v2=off nopti irqpoll nr_cpus=1 reset_devices cgroup_disable=memory mce=off numa=off udev.children-max=2 panic=10 rootflags=nofail acpi_no_memhotplug transparent_hugepage=never nokaslr novmcoredd hest_disable disable_cpu_apicid=0' --initrd=/boot/initramfs-4.18.0-348.7.1.el8_5.x86_64kdump.img /boot/vmlinuz-4.18.0-348.7.1.el8_5.x86_64\"}"
                },
                "@version":"1",
                "type":"carey1"
            }
        }
    ]
}
*/
type ElasticsearchInfra interface {
	SearchRequest(indexNames []string, quer *entity.QueryDocs) (*elastic.SearchResult, error)
	IndicesDeleteRequest(indexNames []string) (*elastic.Response, error)
}

var _ ElasticsearchInfra = new(elasticsearch)

type elasticsearch struct {
	Client *elastic.Client
}

func NewElasticsearch(conf common.AppConfig) (*elasticsearch, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(conf.Elasticsearch.Address...),
		elastic.SetBasicAuth(conf.Elasticsearch.Username, conf.Elasticsearch.Password),
	)

	if err != nil {
		return nil, err
	}
	return &elasticsearch{client}, nil
}

func (es *elasticsearch) SearchRequest(indexNames []string, quer *entity.QueryDocs) (*elastic.SearchResult, error) {
	query := elastic.NewMatchAllQuery()
	res, err := es.Client.Search().Index(indexNames...).From(quer.From).Size(quer.Size).SortBy(quer.Sort...).Query(query).Do(context.Background())

	if err != nil {
		return nil, err
	}
	if res == nil {
		err = errors.New("got results = nil")
		return nil, err
	}
	if res.Hits == nil {
		err = errors.New("got SearchResult.Hits = nil")
		return nil, err
	}

	return res, nil
}


func (es *elasticsearch)IndicesDeleteRequest(indexNames []string) (*elastic.Response, error){
	return 
}