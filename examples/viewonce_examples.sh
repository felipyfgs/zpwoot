#!/bin/bash

# ViewOnce Examples
# This script demonstrates how to send ViewOnce messages using the new parameter approach

API_URL="http://localhost:8080"
API_KEY="your-api-key-here"
SESSION_ID="my-session"
PHONE="5511999999999"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== ViewOnce Message Examples ===${NC}\n"

# Example 1: Send ViewOnce Image
echo -e "${GREEN}Example 1: Send ViewOnce Image${NC}"
curl -X POST "${API_URL}/sessions/${SESSION_ID}/send/message/image" \
  -H "Authorization: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "'"${PHONE}"'",
    "file": "https://picsum.photos/800/600",
    "caption": "Esta imagem desaparecerá após visualização! 🔒",
    "viewOnce": true
  }'
echo -e "\n"

# Example 2: Send ViewOnce Video
echo -e "${GREEN}Example 2: Send ViewOnce Video${NC}"
curl -X POST "${API_URL}/sessions/${SESSION_ID}/send/message/video" \
  -H "Authorization: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "'"${PHONE}"'",
    "file": "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_1mb.mp4",
    "caption": "Vídeo confidencial - ViewOnce 🎥",
    "viewOnce": true
  }'
echo -e "\n"

# Example 3: Send ViewOnce Audio
echo -e "${GREEN}Example 3: Send ViewOnce Audio${NC}"
curl -X POST "${API_URL}/sessions/${SESSION_ID}/send/message/audio" \
  -H "Authorization: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "'"${PHONE}"'",
    "file": "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3",
    "viewOnce": true
  }'
echo -e "\n"

# Example 4: Send ViewOnce Image with Reply (ContextInfo)
echo -e "${GREEN}Example 4: Send ViewOnce Image with Reply${NC}"
curl -X POST "${API_URL}/sessions/${SESSION_ID}/send/message/image" \
  -H "Authorization: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "'"${PHONE}"'",
    "file": "https://picsum.photos/800/600",
    "caption": "Respondendo com uma imagem ViewOnce",
    "viewOnce": true,
    "contextInfo": {
      "stanzaId": "3EB0A9253FA64269E11C9D"
    }
  }'
echo -e "\n"

# Example 5: Send ViewOnce Image with Base64
echo -e "${GREEN}Example 5: Send ViewOnce Image with Base64${NC}"
curl -X POST "${API_URL}/sessions/${SESSION_ID}/send/message/image" \
  -H "Authorization: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "'"${PHONE}"'",
    "file": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
    "caption": "Imagem Base64 ViewOnce",
    "viewOnce": true
  }'
echo -e "\n"

# Example 6: Regular Image (without ViewOnce)
echo -e "${GREEN}Example 6: Regular Image (without ViewOnce for comparison)${NC}"
curl -X POST "${API_URL}/sessions/${SESSION_ID}/send/message/image" \
  -H "Authorization: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "'"${PHONE}"'",
    "file": "https://picsum.photos/800/600",
    "caption": "Esta é uma imagem normal (não ViewOnce)",
    "viewOnce": false
  }'
echo -e "\n"

echo -e "${BLUE}=== Examples completed ===${NC}"
echo -e "${YELLOW}Note: Replace API_KEY, SESSION_ID, and PHONE with your actual values${NC}"
echo -e "${GREEN}All examples use the new viewOnce parameter approach!${NC}"

