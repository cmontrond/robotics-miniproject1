package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	g "gobot.io/x/gobot/platforms/dexter/gopigo3"
	"gobot.io/x/gobot/platforms/raspi"
	"time"
)

const (
	LIGHT_IN_REACH  = 1000
	LIGHT_TOO_CLOSE = 3000
)

func stopRobot(gopigo3 *g.Driver) {
	err := gopigo3.SetMotorDps(g.MOTOR_LEFT, 0)
	if err != nil {
		fmt.Errorf("Error stopping left wheel %+v", err)
	}
	err = gopigo3.SetMotorDps(g.MOTOR_RIGHT, 0)
	if err != nil {
		fmt.Errorf("Error stopping right wheel %+v", err)
	}
}

func robotRunLoop(gopigo3 *g.Driver, leftLightSensor *aio.GroveLightSensorDriver, rightLightSensor *aio.GroveLightSensorDriver) {

	//lightFound := false
	//turnedStraight := false

	for {

		_ = gopigo3.SetMotorDps(g.MOTOR_LEFT, 60)
		_ = gopigo3.SetMotorDps(g.MOTOR_RIGHT, 60)

		time.Sleep(time.Second * 3)

		stopRobot(gopigo3)

		time.Sleep(time.Second * 3)

		//leftLightSensorVal, err := leftLightSensor.Read()
		//
		//if err != nil {
		//	fmt.Errorf("Error reading sensor %+v", err)
		//}
		//
		//fmt.Println("Left Light Value is ", leftLightSensorVal)
		//
		//rightLightSensorVal, err := rightLightSensor.Read()
		//
		//if err != nil {
		//	fmt.Errorf("Error reading sensor %+v", err)
		//}
		//
		//fmt.Println("Right Light Value is ", rightLightSensorVal)
		//
		//if rightLightSensorVal > leftLightSensorVal && rightLightSensorVal >= LIGHT_IN_REACH {
		//
		//} else if leftLightSensorVal > rightLightSensorVal && leftLightSensorVal >= LIGHT_IN_REACH {
		//
		//}

		//if lightFound == false && leftLightSensorVal >= LIGHT_IN_REACH {
		//
		//	fmt.Printf("Turn wheel 1: left sensor")
		//
		//	lightFound = true
		//
		//	if turnedStraight == false {
		//		fmt.Printf("Turn wheel 2: left sensor")
		//		if leftLightSensorVal <= (LIGHT_IN_REACH + 1000) {
		//			//turnedStraight = true
		//			gopigo3.SetMotorPosition(g.MOTOR_LEFT, -90)
		//		}
		//	} else {
		//		fmt.Printf("Turn wheel 3: left sensor")
		//		_ = gopigo3.SetMotorDps(g.MOTOR_LEFT, 150)
		//		_ = gopigo3.SetMotorDps(g.MOTOR_RIGHT, 150)
		//	}
		//}
		//
		//fmt.Println("Right Light Value is ", rightLightSensorVal)
		//
		//if lightFound == false && rightLightSensorVal >= LIGHT_IN_REACH {
		//	fmt.Printf("Turn wheel 1: right sensor")
		//	lightFound = true
		//
		//	if turnedStraight == false {
		//		fmt.Printf("Turn wheel 2: right sensor")
		//		if rightLightSensorVal <= (LIGHT_IN_REACH + 1000) {
		//			//turnedStraight = true
		//			gopigo3.SetMotorPosition(g.MOTOR_RIGHT, 90)
		//		}
		//	} else {
		//		fmt.Printf("Turn wheel 3: right sensor")
		//		_ = gopigo3.SetMotorDps(g.MOTOR_LEFT, 150)
		//		_ = gopigo3.SetMotorDps(g.MOTOR_RIGHT, 150)
		//	}
		//
		//}

		//time.Sleep(time.Second)
		//lightFound = false
	}
}

func main() {

	raspiAdaptor := raspi.NewAdaptor()
	gopigo3 := g.NewDriver(raspiAdaptor)

	leftLightSensor := aio.NewGroveLightSensorDriver(gopigo3, "AD_2_1")
	rightLightSensor := aio.NewGroveLightSensorDriver(gopigo3, "AD_1_1")

	mainRobotFunc := func() {
		robotRunLoop(gopigo3, leftLightSensor, rightLightSensor)
	}

	robot := gobot.NewRobot("miniProject",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{gopigo3, leftLightSensor, rightLightSensor},
		mainRobotFunc,
	)

	err := robot.Start()

	if err != nil {
		fmt.Errorf("Error starting the Robot %+v", err)
	}
}
