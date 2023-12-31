FROM node:20.2.0-alpine

RUN apk update && apk add bash
RUN npm install --global pnpm

ARG BASE_URL
ARG GOOGLE_CLIENT_ID
ARG HOSTED=false
ARG REST_API_BASE_URL
ARG ROBOTS_TXT_ALLOW=false
ARG SENTRY_AUTH_TOKEN
ARG SENTRY_DSN
ARG SENTRY_ENVIRONMENT
ARG SENTRY_ORG
ARG SENTRY_PROJECT
ARG STRIPE_PUBLISHABLE_KEY

WORKDIR /poeticmetric-frontend-builder

# copy only package definition files
COPY package.json .
COPY pnpm-lock.yaml .

# install dependencies
RUN pnpm install

# copy the rest of the files
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
COPY .eslintrc.js .
COPY .stylelintrc.js .
COPY next-sitemap.config.js .
COPY next.config.js .
COPY sentry.client.config.js .
COPY sentry.server.config.js .
COPY tsconfig.json .

ENV NEXT_TELEMETRY_DISABLED=1

ENV NEXT_PUBLIC_BASE_URL=${BASE_URL}
ENV NEXT_PUBLIC_GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID}
ENV NEXT_PUBLIC_HOSTED=${HOSTED}
ENV NEXT_PUBLIC_REST_API_BASE_URL=${REST_API_BASE_URL}
ENV NEXT_PUBLIC_SENTRY_DSN=${SENTRY_DSN}
ENV NEXT_PUBLIC_SENTRY_ENVIRONMENT=${SENTRY_ENVIRONMENT}
ENV NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=${STRIPE_PUBLISHABLE_KEY}
ENV ROBOTS_TXT_ALLOW=${ROBOTS_TXT_ALLOW}
ENV SENTRY_AUTH_TOKEN=${SENTRY_AUTH_TOKEN}
ENV SENTRY_ORG=${SENTRY_ORG}
ENV SENTRY_PROJECT=${SENTRY_PROJECT}

COPY scripts/build.sh .

RUN ./build.sh

FROM nginx:1.23-alpine

RUN apk update && apk add bash

WORKDIR /poeticmetric-frontend

COPY --from=0 /poeticmetric-frontend-builder/out /usr/share/nginx/html
COPY scripts/generate-config.sh /docker-entrypoint.d/40-generate-config.sh
COPY etc/nginx.nginx /etc/nginx/conf.d/default.conf
