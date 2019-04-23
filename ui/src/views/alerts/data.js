var alertItems =  {
    apm: [
    // {
    //   name: 'apm.apdex.count',
    //   label: '综合健康指数Apdex',
    //   compare: 3,
    //   unit: '',
    //   duration: 1,
    //   keys: '',
    //   value: 0.8
    // },
    { 
        name: 'apm.http_code.ratio',
        label: '错误HTTP CODE比率',
        compare: 1,
        unit: '%',
        duration: 1,
        keys: '',
        value: 10,
        help: '指定的http code占所有请求的比例'
    },
    {
        name: 'apm.http_code.count',
        label: '错误HTTP CODE次数',
        compare: 1,
        unit: '次',
        duration: 1,
        keys: '',
        value: 10,
        help: '制定的http code发生次数'
    },
    { 
        name: 'apm.api_error.ratio',
        label: '接口错误率',
        compare: 1,
        unit: '%',
        duration: 1,
        keys: '',
        value: 10
    },
    {
        name: 'apm.sql_error.ratio',
        label: 'sql错误率',
        compare: 1,
        unit: '%',
        duration: 1,
        keys: '',
        value: 10
    },
    {
        name: 'apm.api.duration',
        label: '接口平均耗时',
        compare: 1,
        unit: 'ms',
        duration: 1,
        keys: '',
        value: 10000
    },
        {
        name: 'apm.sql.duration',
        label: 'sql平均耗时',
        compare: 1,
        unit: 'ms',
        duration: 1,
        keys: '',
        value: 10000
    },
    {
        name: 'apm.jvm_fullgc.count',
        label: 'JVMFullGC报警',
        compare: 1,
        unit: '次',
        duration: 1,
        keys: '',
        value: 2
    },
    {
        name: 'apm.api.count',
        label: '接口访问次数',
        compare: 1,
        unit: '次',
        duration: 1,
        keys: '',
        value: 3000
    }
    ],
    system: [
    {
        name: 'system.cpu_used.ratio',
        label: 'cpu使用率',
        compare: 1,
        unit: '%',
        duration: 1,
        keys: '',
        value: 80
    },
    {
        name: 'system.load.count',
        label: '系统Load',
        compare: 1,
        unit: '',
        duration: 1,
        keys: '',
        value: 4
    },
    {
        name: 'system.mem_used.ratio',
        label: '内存使用率',
        compare: 1,
        unit: '%',
        duration: 1,
        keys: '',
        value: 90
    },
    {
        name: 'system.disk_used.ratio',
        label: '硬盘使用率',
        compare: 1,
        unit: '%',
        duration: 1,
        keys: '',
        value: 80
    },
    {
        name: 'system.syn_recv.count',
        label: 'sync_recv数',
        compare: 1,
        unit: '个',
        duration: 1,
        keys: '',
        value: 10000
    },
        {
        name: 'system.time_wait.count',
        label: 'time_wait数',
        compare: 1,
        unit: '个',
        duration: 1,
        keys: '',
        value: 10000
    },
    {
        name: 'system.diskio.ratio',
        label: 'diskio利用率',
        compare: 1,
        unit: '%',
        duration: 1,
        keys: '',
        value: 90
    },
    {
        name: 'system.ifstat_out.speed',
        label: '网络out速度',
        compare: 1,
        unit: 'MB/S',
        duration: 1,
        keys: '',
        value: 100
    },
    {
        name: 'system.close_wait.count',
        label: 'close_wait数',
        compare: 1,
        unit: '个',
        duration: 1,
        keys: '',
        value: 5000
    },
    {
        name: 'system.ifstat_in.speed',
        label: '网络in速度',
        compare: 1,
        unit: 'MB/S',
        duration: 1,
        keys: '',
        value: 100
    },
    {
        name: 'system.estab.count',
        label: '建立长链接数',
        compare: 1,
        unit: '个',
        duration: 1,
        keys: '',
        value: 5000
    }
    ]
};

export default alertItems;
   