#define SPEED_IN_1 (2)
#define SPEED_IN_2 (4)
#define SPEED_EN (5)
#define DIRECTION_IN_1 (A2)
#define DIRECTION_IN_2 (A3)
#define DIRECTION_EN (3)

#define FORARD_BIT   (1) // '0001'
#define BACKWARD_BIT (2) // '0010'
#define LEFT_BIT     (4) // '0100'
#define RIGHT_BIT    (8) // '1000'

void setup()
{
    pinMode(SPEED_IN_1, OUTPUT);
    pinMode(SPEED_IN_2, OUTPUT);
    pinMode(SPEED_EN, OUTPUT);
    pinMode(DIRECTION_IN_1, OUTPUT);
    pinMode(DIRECTION_IN_2, OUTPUT);
    pinMode(DIRECTION_EN, OUTPUT);

    Serial.begin(115200);
}

void drive(byte direction)
{
    if ((direction & FORARD_BIT) && (direction & BACKWARD_BIT)) {
        direction -= BACKWARD_BIT;
    }
    if ((direction & LEFT_BIT) && (direction & RIGHT_BIT)) {
        direction -= RIGHT_BIT;
    }

    digitalWrite(SPEED_IN_1, (direction & FORARD_BIT) ? HIGH : LOW);
    digitalWrite(SPEED_IN_2, (direction & BACKWARD_BIT) ? HIGH : LOW);

    if (direction & FORARD_BIT || direction & BACKWARD_BIT) {
      analogWrite(SPEED_EN, 255);
    } else {
      analogWrite(SPEED_EN, 0);
    }

    digitalWrite(DIRECTION_IN_1, (direction & RIGHT_BIT) ? HIGH : LOW);
    digitalWrite(DIRECTION_IN_2, (direction & LEFT_BIT) ? HIGH : LOW);
    
    if (direction & RIGHT_BIT || direction & LEFT_BIT) {
      analogWrite(DIRECTION_EN, 255);
    } else {
      analogWrite(DIRECTION_EN, 0);
    }
}

byte in;

void loop()
{
  if (Serial.available()) {
    in = Serial.read();
    drive(in);
    Serial.println(in, DEC);
  }
}

