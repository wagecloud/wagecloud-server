-- CreateEnum
CREATE TYPE "PaymentType" AS ENUM ('VNPay');

-- CreateTable
CREATE TABLE "account" (
    "id" BIGINT NOT NULL,
    "email" TEXT NOT NULL,
    "name" TEXT,
    "password" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "account_pkey" PRIMARY KEY ("id")
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
CREATE UNIQUE INDEX "account_email_key" ON "account"("email");

-- CreateIndex
CREATE UNIQUE INDEX "vm_network_id_key" ON "vm"("network_id");

-- CreateIndex
CREATE UNIQUE INDEX "network_private_ip_key" ON "network"("private_ip");

-- AddForeignKey
ALTER TABLE "vm" ADD CONSTRAINT "vm_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "account"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vm" ADD CONSTRAINT "vm_network_id_fkey" FOREIGN KEY ("network_id") REFERENCES "network"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vm" ADD CONSTRAINT "vm_os_id_fkey" FOREIGN KEY ("os_id") REFERENCES "os"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "vm" ADD CONSTRAINT "vm_arch_id_fkey" FOREIGN KEY ("arch_id") REFERENCES "arch"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

