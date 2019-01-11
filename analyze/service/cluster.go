package service

// // Cluster ...
// type Cluster struct {
// }

// type eventDelegate struct {
// }

// // NotifyJoin is invoked when a node is detected to have joined.
// // The Node argument must not be modified.
// func (e *eventDelegate) NotifyJoin(n *memberlist.Node) {
// 	gAnalyze.hash.Add(n.Name)
// 	g.L.Info("NotifyJoin", zap.String("name", n.Name))
// }

// // NotifyLeave is invoked when a node is detected to have left.
// // The Node argument must not be modified.
// func (e *eventDelegate) NotifyLeave(n *memberlist.Node) {
// 	gAnalyze.hash.Remove(n.Name)
// 	g.L.Info("NotifyLeave", zap.String("name", n.Name))
// }

// func (e *eventDelegate) NotifyUpdate(n *memberlist.Node) {
// 	g.L.Info("NotifyUpdate", zap.String("name", n.Name))
// }

// // NewCluster ...
// func NewCluster() *Cluster {
// 	return &Cluster{}
// }

// // Start ...
// func (cluster *Cluster) Start() error {

// 	config := memberlist.DefaultLocalConfig()

// 	host, err := os.Hostname()
// 	if err != nil {
// 		g.L.Fatal("get host name", zap.String("error", err.Error()))
// 	}

// 	if misc.Conf.Cluster.HostUseTime {
// 		config.Name = host + time.Now().UTC().String()
// 	} else {
// 		config.Name = host
// 	}

// 	misc.Conf.Cluster.Name = config.Name
// 	gAnalyze.hash.Add(config.Name)

// 	config.BindAddr = misc.Conf.Cluster.Addr
// 	config.BindPort = misc.Conf.Cluster.Port

// 	config.AdvertiseAddr = misc.Conf.Cluster.Addr
// 	config.AdvertisePort = misc.Conf.Cluster.Port
// 	config.Events = &eventDelegate{}

// 	list, err := memberlist.Create(config)
// 	if err != nil {
// 		g.L.Panic("Cluster Start", zap.Error(err))
// 	}

// 	_, err = list.Join(misc.Conf.Cluster.Seeds)
// 	if err != nil {
// 		g.L.Panic("Cluster Join", zap.Error(err))
// 	}

// 	return nil
// }

// // Close ...
// func (cluster *Cluster) Close() error {
// 	return nil
// }
