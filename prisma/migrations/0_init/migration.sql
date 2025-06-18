-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "account";

-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "instance";

-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "os";

-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "payment";

-- CreateEnum
CREATE TYPE "account"."type" AS ENUM ('ACCOUNT_TYPE_ADMIN', 'ACCOUNT_TYPE_USER');

-- CreateEnum
CREATE TYPE "instance"."log_type" AS ENUM ('LOG_TYPE_UNKNOWN', 'LOG_TYPE_INFO', 'LOG_TYPE_WARNING', 'LOG_TYPE_ERROR');

-- CreateEnum
CREATE TYPE "payment"."method" AS ENUM ('PAYMENT_METHOD_UNKNOWN', 'PAYMENT_METHOD_VNPAY', 'PAYMENT_METHOD_MOMO');

-- CreateEnum
CREATE TYPE "payment"."status" AS ENUM ('PAYMENT_STATUS_UNKNOWN', 'PAYMENT_STATUS_PENDING', 'PAYMENT_STATUS_SUCCESS', 'PAYMENT_STATUS_CANCELED', 'PAYMENT_STATUS_FAILED');

-- CreateTable
CREATE TABLE "account"."base" (
    "id" BIGSERIAL NOT NULL,
    "type" "account"."type" NOT NULL,
    "username" TEXT NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."user" (
    "id" BIGINT NOT NULL,
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "email" TEXT,
    "phone" TEXT,
    "company" VARCHAR(255),
    "address" VARCHAR(255),

    CONSTRAINT "user_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "instance"."base" (
    "id" TEXT NOT NULL,
    "account_id" BIGINT NOT NULL,
    "os_id" TEXT NOT NULL,
    "arch_id" TEXT NOT NULL,
    "region_id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "cpu" INTEGER NOT NULL,
    "ram" INTEGER NOT NULL,
    "storage" INTEGER NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "instance"."network" (
    "id" BIGSERIAL NOT NULL,
    "instance_id" TEXT NOT NULL,
    "private_ip" TEXT NOT NULL,
    "mac_address" TEXT NOT NULL,
    "public_ip" TEXT,

    CONSTRAINT "network_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "instance"."domain" (
    "id" BIGSERIAL NOT NULL,
    "network_id" BIGINT NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "domain_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "instance"."log" (
    "id" BIGSERIAL NOT NULL,
    "instance_id" TEXT NOT NULL,
    "type" "instance"."log_type" NOT NULL,
    "title" TEXT NOT NULL,
    "description" TEXT,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "log_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "instance"."region" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "region_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "os"."base" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "os"."arch" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "arch_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "payment"."item" (
    "id" BIGSERIAL NOT NULL,
    "payment_id" BIGINT NOT NULL,
    "name" TEXT NOT NULL,
    "price" BIGINT NOT NULL,

    CONSTRAINT "item_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "payment"."base" (
    "id" BIGSERIAL NOT NULL,
    "account_id" BIGINT NOT NULL,
    "method" "payment"."method" NOT NULL,
    "status" "payment"."status" NOT NULL,
    "total" BIGINT NOT NULL,
    "date_created" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "payment"."vnpay" (
    "id" BIGINT NOT NULL,
    "vnp_TxnRef" TEXT NOT NULL,
    "vnp_OrderInfo" TEXT NOT NULL,
    "vnp_TransactionNo" TEXT NOT NULL,
    "vnp_TransactionDate" TEXT NOT NULL,
    "vnp_CreateDate" TEXT NOT NULL,
    "vnp_IpAddr" TEXT NOT NULL,

    CONSTRAINT "vnpay_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "base_username_key" ON "account"."base"("username");

-- CreateIndex
CREATE UNIQUE INDEX "user_email_key" ON "account"."user"("email");

-- CreateIndex
CREATE UNIQUE INDEX "user_phone_key" ON "account"."user"("phone");

-- CreateIndex
CREATE UNIQUE INDEX "network_instance_id_key" ON "instance"."network"("instance_id");

-- CreateIndex
CREATE UNIQUE INDEX "domain_name_key" ON "instance"."domain"("name");

-- AddForeignKey
ALTER TABLE "account"."user" ADD CONSTRAINT "user_id_fkey" FOREIGN KEY ("id") REFERENCES "account"."base"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "instance"."base" ADD CONSTRAINT "base_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "account"."user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "instance"."base" ADD CONSTRAINT "base_os_id_fkey" FOREIGN KEY ("os_id") REFERENCES "os"."base"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "instance"."base" ADD CONSTRAINT "base_arch_id_fkey" FOREIGN KEY ("arch_id") REFERENCES "os"."arch"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "instance"."base" ADD CONSTRAINT "base_region_id_fkey" FOREIGN KEY ("region_id") REFERENCES "instance"."region"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "instance"."network" ADD CONSTRAINT "network_instance_id_fkey" FOREIGN KEY ("instance_id") REFERENCES "instance"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "instance"."domain" ADD CONSTRAINT "domain_network_id_fkey" FOREIGN KEY ("network_id") REFERENCES "instance"."network"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "instance"."log" ADD CONSTRAINT "log_instance_id_fkey" FOREIGN KEY ("instance_id") REFERENCES "instance"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "payment"."item" ADD CONSTRAINT "item_payment_id_fkey" FOREIGN KEY ("payment_id") REFERENCES "payment"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "payment"."base" ADD CONSTRAINT "base_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "account"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "payment"."vnpay" ADD CONSTRAINT "vnpay_id_fkey" FOREIGN KEY ("id") REFERENCES "payment"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

