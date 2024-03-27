import json
from datetime import datetime
from garmin_fit_sdk import Decoder, Stream
import sys
import os

# pip3 install garmin-fit-sdk

class Point:

    def __init__(self, ts, lat, lng, distance, speed, altitude):
        self.ts = int(ts)
        self.lat = lat
        self.lng = lng
        self.distance = distance
        self.speed = speed
        self.altitude = altitude


def main(f: str):
    if not os.path.exists(f):
        print(f'file: {f} not exist, exit')
        return
    stream = Stream.from_file(f)
    decoder = Decoder(stream)
    messages, errors = decoder.read()

    points = messages['record_mesgs']

    results: list[Point] = []
    for point in points:
        ts = datetime.timestamp(point['timestamp'])
        lat = point['position_lat'] * (180 / 2 ** 31)
        lng = point['position_long'] * (180 / 2 ** 31)
        distance = point['distance']
        speed = point['enhanced_speed']
        altitude = point['enhanced_altitude']

        p = Point(ts, lat, lng, distance, speed, altitude)
        results.append(p)

    with open('result.json', 'w') as f:
        json.dump(results, f, default=vars)


if __name__ == '__main__':
    args = sys.argv
    if len(args) == 1:
        print('no fit file param')
        sys.exit(0)
    main(args[1])
