import json
import os
import sys
from datetime import datetime

from garmin_fit_sdk import Decoder, Stream


# pip3 install garmin-fit-sdk

class Point:
    distance = 0
    speed = 0
    cadence = 0
    altitude = 0
    hr = 0

    def __init__(self, ts, lat, lng):
        self.ts = int(ts)
        self.lat = lat
        self.lng = lng


def main(f: str, result: str):
    if not os.path.exists(f):
        print(f'file: {f} not exist, exit')
        return
    stream = Stream.from_file(f)
    decoder = Decoder(stream)
    messages, errors = decoder.read()

    points = messages['record_mesgs']

    results: list[Point] = []
    for point in points:
        timestamp = point.get('timestamp')
        if not timestamp:
            break
        ts = datetime.timestamp(timestamp)
        position_lat = point.get('position_lat')
        if position_lat is None:
            continue
        lat = point['position_lat'] * (180 / 2 ** 31)
        lng = point['position_long'] * (180 / 2 ** 31)

        p = Point(ts, lat, lng)

        p.hr = point.get('heart_rate', 0)
        p.altitude = point.get('enhanced_altitude', 0)
        p.distance = point.get('distance', 0)
        p.speed = point.get('enhanced_speed', 0)
        p.cadence = point.get('cadence', 0)

        results.append(p)

    with open(result, 'w') as f:
        json.dump(results, f, default=vars)


if __name__ == '__main__':
    args = sys.argv
    if len(args) < 3:
        # python3 fit_parser.py xxx.fit xxx.json
        print('no fit file param')
        sys.exit(0)
    main(args[1], args[2])
