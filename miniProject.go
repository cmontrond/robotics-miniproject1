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
	LIGHT_IN_REACH  = 2500
	LIGHT_TOO_CLOSE = 3500
	TURN_POSITION   = -10
	FORWARD_DPS     = -50
)

func stop(gopigo3 *g.Driver) {
	err := gopigo3.SetMotorDps(g.MOTOR_LEFT, 0)
	if err != nil {
		fmt.Errorf("Error stopping left wheel %+v", err)
	}
	err = gopigo3.SetMotorDps(g.MOTOR_RIGHT, 0)
	if err != nil {
		fmt.Errorf("Error stopping right wheel %+v", err)
	}
}

func turnRight(gopigo3 *g.Driver) {
	err := gopigo3.SetMotorPosition(g.MOTOR_LEFT, TURN_POSITION)
	if err != nil {
		fmt.Errorf("Error turning right wheel %+v", err)
	}
}

func turnLeft(gopigo3 *g.Driver) {
	err := gopigo3.SetMotorPosition(g.MOTOR_RIGHT, TURN_POSITION)
	if err != nil {
		fmt.Errorf("Error turning left wheel %+v", err)
	}
}

func moveForward(gopigo3 *g.Driver) {
	err := gopigo3.SetMotorDps(g.MOTOR_LEFT, FORWARD_DPS)
	if err != nil {
		fmt.Errorf("Error moving left wheel %+v", err)
	}
	err = gopigo3.SetMotorDps(g.MOTOR_RIGHT, FORWARD_DPS)
	if err != nil {
		fmt.Errorf("Error moving right wheel %+v", err)
	}
}

func blinkLED(gopigo3 *g.Driver) {
	err := gopigo3.SetLED(g.LED_EYE_RIGHT, 0x00, 0x00, 0xFF)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second)

	err = gopigo3.SetLED(g.LED_EYE_RIGHT, 0x00, 0x00, 0x00)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second)
}

func robotRunLoop(gopigo3 *g.Driver, leftLightSensor *aio.GroveLightSensorDriver, rightLightSensor *aio.GroveLightSensorDriver) {

	robotStopped := false

	for {

		// Read value from the left light sensor
		leftLightSensorVal, err := leftLightSensor.Read()

		if err != nil {
			fmt.Errorf("Error reading sensor %+v", err)
		}

		// read value from the right light sensor
		rightLightSensorVal, err := rightLightSensor.Read()

		if err != nil {
			fmt.Errorf("Error reading sensor %+v", err)
		}

		fmt.Println("Right Light Value is ", rightLightSensorVal)
		fmt.Println("Left Light Value is ", leftLightSensorVal)

		// Stop the Robot if too close to the light
		if rightLightSensorVal >= LIGHT_TOO_CLOSE || leftLightSensorVal >= LIGHT_TOO_CLOSE {
			stop(gopigo3)
			robotStopped = true
			blinkLED(gopigo3)
		}

		if robotStopped == false {

			rightLeftDifference := rightLightSensorVal - leftLightSensorVal
			leftRightDifference := leftLightSensorVal - rightLightSensorVal

			// If the light comes from the right, turn right and move forward
			if (rightLightSensorVal > leftLightSensorVal) && (rightLightSensorVal >= LIGHT_IN_REACH) && (rightLeftDifference >= 150) {

				turnRight(gopigo3)
				time.Sleep(time.Second)
				moveForward(gopigo3)
				time.Sleep(time.Second)

				// If the light comes from the left, turn left and move forward
			} else if (leftLightSensorVal > rightLightSensorVal) && (leftLightSensorVal >= LIGHT_IN_REACH) && (leftRightDifference >= 150) {

				turnLeft(gopigo3)
				time.Sleep(time.Second)
				moveForward(gopigo3)
				time.Sleep(time.Second)

			} else if (rightLightSensorVal >= LIGHT_IN_REACH) && (leftLightSensorVal >= LIGHT_IN_REACH) {
				moveForward(gopigo3)
				time.Sleep(time.Second)
			}
		}
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
