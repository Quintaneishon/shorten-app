import requests
import json
import time
import os

url = os.environ.get('URL', 'http://localhost/shorten')

# list of JSON objects with different data
numbers = list(range(100))

start_time = time.time()
responses = []

# send 100 POST requests in 2 seconds
for i in range(100):
    # randomly select a JSON object from the list
    headers = {'Content-type': 'application/json'}
    response = requests.post(url, data=json.dumps({"long_url":f"http://ejemplo.com/{numbers[i]}"}), headers=headers)
    responses.append(response)
    time.sleep(0.02) # sleep for 20ms to achieve 50 requests per second

end_time = time.time()
total_time = end_time - start_time

# calculate the average response time
avg_response_time = sum([r.elapsed.total_seconds() for r in responses])/len(responses) * 1000

# calculate the percentage of requests that took less than 10ms
less_than_10ms = sum([1 for r in responses if r.elapsed.total_seconds() < 0.01])/len(responses)*100

print(f'Average response time: {avg_response_time:.2f} ms')
print(f'Percentage of requests less than 10ms: {less_than_10ms:.2f}%')
print(f'Total time: {total_time:.2f} seconds')
