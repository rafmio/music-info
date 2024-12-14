curl -X POST \
  http://localhost:8080/v1/songs \
  -H 'Content-Type: application/json' \
  -d '{
  "group": "Dire Straits",
  "song": "Why Worry"
}'
