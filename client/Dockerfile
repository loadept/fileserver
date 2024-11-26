FROM node:22-alpine3.19 AS build

WORKDIR /app

ARG SECRET_API_URL
ENV SECRET_API_URL=${SECRET_API_URL}

COPY package*.json ./

RUN npm ci

COPY . .

RUN npm run build

FROM node:22-alpine3.19

WORKDIR /app

COPY --from=build /app/dist .
COPY --from=build /app/node_modules node_modules

ENTRYPOINT [ "node", "server/entry.mjs" ]
