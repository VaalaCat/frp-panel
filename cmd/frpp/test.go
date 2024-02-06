package main

// func serverThings() {
// 	time.Sleep(5 * time.Second)
// 	p := pool.New().WithMaxGoroutines(500)
// 	for i := 0; i < 100000; i++ {
// 		cnt := i
// 		p.Go(func() {
// 			resp := rpc.CallClient(context.Background(), "test", &pb.ServerMessage{
// 				Event: pb.Event_EVENT_DATA,
// 				Data:  []byte(fmt.Sprint(cnt)),
// 			})
// 			if string(resp.Data) != fmt.Sprint(cnt) {
// 				logrus.Panicf("resp: %+v", resp)
// 			}
// 		})
// 	}
// 	p.Wait()
// 	logrus.Infof("finish server things")
// }
