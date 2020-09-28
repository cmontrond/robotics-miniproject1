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
	LIGHT_IN_REACH  = 2000
	LIGHT_TOO_CLOSE = 100
)

func robotRunLoop(gopigo3 *g.Driver, leftLightSensor *aio.GroveLightSensorDriver, rightLightSensor *aio.GroveLightSensorDriver) {

	lightFound := false
	turnedStraight := false

	for {

		leftLightSensorVal, err := leftLightSensor.Read()

		if err != nil {
			fmt.Errorf("Error reading sensor %+v", err)
		}

		fmt.Println("Left Light Value is ", leftLightSensorVal)

		if lightFound == false && leftLightSensorVal >= LIGHT_IN_REACH {
			lightFound = true

			if turnedStraight == false {
				gopigo3.SetMotorPosition(g.MOTOR_LEFT, -90)
				if leftLightSensorVal >= (LIGHT_IN_REACH + 1000) {
					turnedStraight = true
				}
			} else {
				_ = gopigo3.SetMotorDps(g.MOTOR_LEFT, -150)
				_ = gopigo3.SetMotorDps(g.MOTOR_RIGHT, -150)
			}
		}

		rightLightSensorVal, err := leftLightSensor.Read()

		if err != nil {
			fmt.Errorf("Error reading sensor %+v", err)
		}

		fmt.Println("Right Light Value is ", rightLightSensorVal)

		if lightFound == false && rightLightSensorVal >= LIGHT_IN_REACH {
			lightFound = true

			if turnedStraight == false {
				gopigo3.SetMotorPosition(g.MOTOR_RIGHT, 90)
				if leftLightSensorVal >= (LIGHT_IN_REACH + 1000) {
					turnedStraight = true
				}
			} else {
				_ = gopigo3.SetMotorDps(g.MOTOR_LEFT, -150)
				_ = gopigo3.SetMotorDps(g.MOTOR_RIGHT, -150)
			}

		}

		time.Sleep(time.Second)
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
