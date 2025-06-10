-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "account";

-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "instance";

-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "os";

-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "payment";

-- CreateEnum
CREATE TYPE "account"."role" AS ENUM ('ROLE_ADMIN', 'ROLE_USER');

-- CreateEnum
CREATE TYPE "payment"."method" AS ENUM ('VNPAY', 'MOMO');

-- CreateEnum
CREATE TYPE "payment"."status" AS ENUM ('PENDING', 'SUCCESS', 'CANCELED', 'FAILED');

-- CreateTable
CREATE TABLE "account"."base" (
    "id" BIGSERIAL NOT NULL,
    "role" "account"."role" NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "username" TEXT NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."user" (
    "id" BIGINT NOT NULL,
    "email" TEXT NOT NULL,

    CONSTRAINT "user_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "instance"."base" (
    "id" TEXT NOT NULL,
    "account_id" BIGINT NOT NULL,
    "network_id" TEXT NOT NULL,
    "os_id" TEXT NOT NULL,
    "arch_id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "cpu" INTEGER NOT NULL,
    "ram" INTEGER NOT NULL,
    "storage" INTEGER NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "instance"."network" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "private_ip" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "network_pkey" PRIMARY KEY ("id")
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
CREATE TABLE "os"."image" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "os_id" TEXT NOT NULL,
    "arch_id" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "image_pkey" PRIMARY KEY ("id")
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
CREATE UNIQUE INDEX "base_network_id_key" ON "instance"."base"("network_id");

-- AddForeignKey
ALTER TABLE "account"."user" ADD CONSTRAINT "user_id_fkey" FOREIGN KEY ("id") REFERENCES "account"."base"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "instance"."base" ADD CONSTRAINT "base_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "account"."user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "instance"."base" ADD CONSTRAINT "base_network_id_fkey" FOREIGN KEY ("network_id") REFERENCES "instance"."network"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "instance"."base" ADD CONSTRAINT "base_os_id_fkey" FOREIGN KEY ("os_id") REFERENCES "os"."base"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "instance"."base" ADD CONSTRAINT "base_arch_id_fkey" FOREIGN KEY ("arch_id") REFERENCES "os"."arch"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "os"."image" ADD CONSTRAINT "image_os_id_fkey" FOREIGN KEY ("os_id") REFERENCES "os"."base"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "os"."image" ADD CONSTRAINT "image_arch_id_fkey" FOREIGN KEY ("arch_id") REFERENCES "os"."arch"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "payment"."item" ADD CONSTRAINT "item_payment_id_fkey" FOREIGN KEY ("payment_id") REFERENCES "payment"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "payment"."base" ADD CONSTRAINT "base_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "account"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "payment"."vnpay" ADD CONSTRAINT "vnpay_id_fkey" FOREIGN KEY ("id") REFERENCES "payment"."base"("id") ON DELETE CASCADE ON UPDATE CASCADE;

