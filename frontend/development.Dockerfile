FROM node:20.2.0-alpine

RUN apk update && apk add bash
RUN npm install --global pnpm

WORKDIR /poeticmetric

# copy only package definition files
COPY package.json .
COPY pnpm-lock.yaml .

# install dependencies
RUN pnpm install

COPY @types @types
COPY blog blog
COPY components components
COPY contexts contexts
COPY docs docs
COPY helpers helpers
COPY hooks hooks
COPY lib lib
COPY pages pages
COPY public public
COPY styles styles
COPY tests tests
COPY next.config.js .
COPY next-sitemap.config.js .
COPY playwright.config.ts .
COPY sentry.client.config.js .
COPY sentry.server.config.js .
COPY tsconfig.json .

EXPOSE 80

CMD ["pnpm", "run", "dev"]
