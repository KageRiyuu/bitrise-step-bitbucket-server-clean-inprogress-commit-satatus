import requests
import json
import os

DOMAIN = os.environ['DOMAIN']
ACCESS_TOKEN = os.environ['ACCESS_TOKEN'] 
COMMIT = os.environ['COMMIT'] 

def getStatus(commitId):
    import json
    api_url = DOMAIN + "/rest/build-status/1.0/commits/" + commitId
    print("Url: " +  api_url)
    
    headers =  {
        "Authorization":"Bearer " + ACCESS_TOKEN
    }
    print("Headers: " + headers["Authorization"])
    
    response = requests.get(api_url, headers=headers)
    print("Response: " + str(response))
    if(response.ok):
        data = json.loads(response.content)
        return data["values"]
    else:
        return []
    
def updateStatus(commitId, status, key, name, url, description):
    import json
    api_url = "https://" +  DOMAIN + "/rest/build-status/1.0/commits/" + commitId
    print("Url: " +  api_url)

    headers =  {
        "Authorization":"Bearer " + ACCESS_TOKEN
    }
    print("Headers: " + headers["Authorization"])
    
    data = {
        "state": status,
        "key": key,
        "name": name,
        "url": url,
        "description": description
    }
    print("Data Sending: " + str(data))

    response = requests.post(api_url, json=data, headers=headers)
    print("Response: " + str(response))
    return response.ok
    

statuses = getStatus(COMMIT)
print("Total Statuses: " + str(len(statuses)))
print("Got response: ")
print("----------------------------------------------------------------")
print(str(statuses))
print("----------------------------------------------------------------")

print("")
print("")
print("----------------------------------------------------------------")
print("Updating....")
index = 0
updatedCount = 0
for status in statuses:
    if(status["state"] == "INPROGRESS"):
        isUpdated = updateStatus(COMMIT, "FAILED", status["key"], status["name"], status["url"], status["description"])
        print("Updated: " + str(isUpdated) + " for index: " + str(index))
        updatedCount = updatedCount + 1
    index = index + 1
if(updatedCount == 0):
    print("No status to update !!!")

print("The end...")