FROM node:21.6-alpine

WORKDIR /vue_app

expose 8080
CMD ["npm", "run", "serve"]
