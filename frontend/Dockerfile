# build stage
FROM node:14.7.0-alpine3.12 AS build-env

WORKDIR /frontend

RUN apk add --no-cache python3 make g++
COPY package.json package-lock.json ./
RUN npm install

COPY . ./
RUN npm run build

# final image
FROM nginx:1.19.1-alpine
COPY --from=build-env /frontend/build /usr/share/nginx/html
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
