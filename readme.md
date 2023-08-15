# README

![go_version](https://img.shields.io/badge/Go-1.20.4-brightgreen?style=for-the-badge)

[![downloads](https://img.shields.io/github/downloads/SydneyOwl/shx8800-ble-connector/total)](https://github.com/SydneyOwl/shx8800-ble-connector/releases?style=for-the-badge)
[![downloads@latest](https://img.shields.io/github/downloads/SydneyOwl/shx8800-ble-connector/latest/total)](https://github.com/SydneyOwl/shx8800-ble-connector/releases/latest?style=for-the-badge)
![](https://img.shields.io/github/v/tag/sydneyowl/shx8800-ble-connector?label=version&style=flat-square?style=for-the-badge)

![IMG_20230815_083605](./md_assets/readme/IMG_20230815_083605.jpg)

## 简介

该软件能在具有蓝牙能力的pc端上在没有写频线的情况下对森海克斯8800进行写频。

本质上该软件是一个中继器。由于windows无法直接为森海克斯8800蓝牙创建虚拟端口，该软件能做到连接森海克斯的蓝牙，在虚拟端口软件的支持下，与官方写频软件进行数据包的交换。（蓝牙写频协议和写频线串口写频协议应该是一致的，不确定，但是经测试可用）

![struc](./md_assets/readme/struc.png)

## 使用方法

**警告：使用前做好备份工作！绝对不要将此软件用于除读写频之外的功能（如系统升级）！对使用此软件产生的一切后果作者不任何责任！**

### 前置工作

请下载任意虚拟端口软件（以HDD Virtual Serial Port Tools为例），创建两个互通的虚拟端口（以COM1\COM2为例）

![image-20230815103331778](./md_assets/readme/image-20230815103331778.png)

### 手动编译

熟悉go的朋友们可以选择从源码编译。

直接go mod tidy然后go build即可。需要gcc支持。

### 直接使用

直接从releases下载对应版本的软件，双击运行。

按照提示，将本软件连接到COM1端口，写频工具连接到COM2即可。

## 许可证

本软件采用Unlicense许可证