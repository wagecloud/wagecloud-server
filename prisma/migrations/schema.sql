-- CreateEnum
CREATE TYPE "OS" AS ENUM ('Ubuntu', 'Debian');

-- CreateEnum
CREATE TYPE "Arch" AS ENUM ('x86_64');

-- CreateEnum
CREATE TYPE "PaymentType" AS ENUM ('VNPay');

-- CreateTable
CREATE TABLE "account" (
    "id" BIGINT NOT NULL,
    "email" TEXT NOT NULL,
    "name" TEXT,
    "password" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "account_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "vm" (
    "id" BIGINT NOT NULL,
    "account_id" BIGINT NOT NULL,
    "name" TEXT NOT NULL,
    "vcpu" INTEGER NOT NULL,
    "ram" INTEGER NOT NULL,
    "os" "OS" NOT NULL,
    "arch" "Arch" NOT NULL,
    "storage" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "vm_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "networking" (
    "id" BIGINT NOT NULL,
    "vm_id" BIGINT NOT NULL,
    "public_ip" TEXT,
    "virtual_network" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "networking_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "account_email_key" ON "account"("email");

-- CreateIndex
CREATE UNIQUE INDEX "networking_vm_id_key" ON "networking"("vm_id");

-- AddForeignKey
ALTER TABLE "vm" ADD CONSTRAINT "vm_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "account"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "networking" ADD CONSTRAINT "networking_vm_id_fkey" FOREIGN KEY ("vm_id") REFERENCES "vm"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

