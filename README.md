# Cloudydocker

Sample app to try out file upload to Cloudinary from Golang. 
Also I experimented dockerizing deploying to [Zeit Now](https://zeit.co/).

Deploy with:
````
$ now -e CLOUDINARY_API_KEY="<api key>" -e CLOUDINARY_API_SECRET="<api secret>" -e CLOUDINARY_CLOUD_NAME="<bucket>" -e PORT=":80"
````

