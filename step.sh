echo "Domain: $DOMAIN"
echo "Commit: $COMMIT"
echo "Access Token: $ACCESS_TOKEN"

python3 CleanBuildStatus.py -DOMAIN=$DOMAIN -COMMIT=$COMMIT -ACCESS_TOKEN=$ACCESS_TOKEN