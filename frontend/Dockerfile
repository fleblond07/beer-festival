FROM node:20-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm ci --ignore-scripts

COPY src/ ./src
COPY public/ ./public
COPY tsconfig.* ./
COPY vite.config.* ./
COPY .env.prod ./.env
COPY index.html ./

RUN npm run build

FROM node:20-alpine

WORKDIR /app

RUN npm install -g serve --ignore-scripts

COPY --from=build /app/dist ./dist

RUN adduser -D beer-festival
USER festival

EXPOSE 1338

CMD ["serve", "-s", "dist", "-l", "1338"]
