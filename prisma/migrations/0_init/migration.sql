-- CreateEnum
CREATE TYPE "PaymentType" AS ENUM ('VNPay');

-- CreateTable
CREATE TABLE "account_base" (
    "id" BIGSERIAL NOT NULL,
    "username" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "account_base_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account_user" (
    "id" BIGINT NOT NULL,

    CONSTRAINT "account_user_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vm" (
    "id" BIGSERIAL NOT NULL,
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

    CONSTRAINT "vm_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "network" (
    "id" TEXT NOT NULL,
    "private_ip" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "network_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "os" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "os_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "arch" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "arch_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "account_base_username_key" ON "account_base"("username");

-- CreateIndex
CREATE UNIQUE INDEX "account_base_email_key" ON "account_base"("email");

-- CreateIndex
CREATE UNIQUE INDEX "vm_network_id_key" ON "vm"("network_id");

-- CreateIndex
CREATE UNIQUE INDEX "network_private_ip_key" ON "network"("private_ip");

-- AddForeignKey
ALTER TABLE "account_user" ADD CONSTRAINT "account_user_id_fkey" FOREIGN KEY ("id") REFERENCES "account_base"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vm" ADD CONSTRAINT "vm_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "account_user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vm" ADD CONSTRAINT "vm_network_id_fkey" FOREIGN KEY ("network_id") REFERENCES "network"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vm" ADD CONSTRAINT "vm_os_id_fkey" FOREIGN KEY ("os_id") REFERENCES "os"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vm" ADD CONSTRAINT "vm_arch_id_fkey" FOREIGN KEY ("arch_id") REFERENCES "arch"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

