FROM node:20-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm ci --ignore-scripts

COPY .env.prod ./.env

COPY . .

RUN npm run build

FROM node:20-alpine

WORKDIR /app

RUN npm install --ignore-scripts -g serve

COPY --from=build /app/dist ./dist

RUN adduser -D beer-festival
USER beer-festival

EXPOSE 1338

CMD ["serve", "-s", "dist", "-l", "1338"]
