"""
eSIM代理API SDK for Python
版本: 1.0.0
文档: https://your-domain.com/api-docs
"""

import hashlib
import hmac
import json
import random
import string
import time
from typing import Dict, List, Optional, Any
from urllib.parse import urlencode

import requests


class EsimSDK:
    """eSIM代理API Python SDK"""

    def __init__(self, config: Dict[str, Any]):
        """
        初始化SDK

        Args:
            config: 配置字典，包含以下键值：
                - apiKey: API密钥
                - apiSecret: API密钥
                - baseURL: API基础URL（可选，默认为 https://api.your-domain.com）
                - timeout: 请求超时时间（可选，默认为30秒）
        """
        self.api_key = config['apiKey']
        self.api_secret = config['apiSecret']
        self.base_url = config.get('baseURL', 'https://api.your-domain.com')
        self.timeout = config.get('timeout', 30)

        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json',
            'User-Agent': 'EsimSDK-Python/1.0.0'
        })

    def generate_signature(self, method: str, path: str, body: Optional[str], timestamp: str, nonce: str) -> str:
        """生成API签名"""
        body_str = body or ''
        sign_string = f"{method}{path}{body_str}{timestamp}{nonce}"

        return hmac.new(
            self.api_secret.encode('utf-8'),
            sign_string.encode('utf-8'),
            hashlib.sha256
        ).hexdigest()

    @staticmethod
    def generate_nonce(length: int = 16) -> str:
        """生成随机字符串"""
        return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

    def request(self, method: str, path: str, data: Optional[Dict] = None) -> Dict[str, Any]:
        """
        发送API请求

        Args:
            method: HTTP方法
            path: API路径
            data: 请求数据（可选）

        Returns:
            API响应数据

        Raises:
            Exception: 请求失败时抛出异常
        """
        timestamp = str(int(time.time()))
        nonce = self.generate_nonce()

        body_str = json.dumps(data, separators=(',', ':')) if data else None
        signature = self.generate_signature(method, path, body_str, timestamp, nonce)

        headers = {
            'x-api-key': self.api_key,
            'x-timestamp': timestamp,
            'x-nonce': nonce,
            'x-signature': signature
        }

        url = f"{self.base_url}{path}"

        try:
            if method.upper() == 'GET':
                response = self.session.get(url, headers=headers, timeout=self.timeout)
            elif method.upper() == 'POST':
                response = self.session.post(url, headers=headers, json=data, timeout=self.timeout)
            elif method.upper() == 'PUT':
                response = self.session.put(url, headers=headers, json=data, timeout=self.timeout)
            elif method.upper() == 'DELETE':
                response = self.session.delete(url, headers=headers, timeout=self.timeout)
            else:
                raise ValueError(f"不支持的HTTP方法: {method}")

            if response.status_code >= 400:
                error_msg = f"API Error: {response.status_code}"
                try:
                    error_data = response.json()
                    error_msg += f" - {error_data.get('message', response.text)}"
                except:
                    error_msg += f" - {response.text}"
                raise Exception(error_msg)

            return response.json()

        except requests.exceptions.RequestException as e:
            raise Exception(f"Network Error: {str(e)}")

    # ========== 产品管理 ==========

    def get_products(self, **params) -> Dict[str, Any]:
        """
        获取产品列表

        Args:
            **params: 查询参数
                - country: 国家代码
                - type: 产品类型
                - limit: 返回数量限制
                - offset: 偏移量

        Returns:
            产品列表响应
        """
        query_string = urlencode({k: v for k, v in params.items() if v is not None})
        path = f"/api/products?{query_string}" if query_string else "/api/products"
        return self.request('GET', path)

    def get_product(self, product_id: str) -> Dict[str, Any]:
        """
        获取产品详情

        Args:
            product_id: 产品ID

        Returns:
            产品详情响应
        """
        return self.request('GET', f"/api/products/{product_id}")

    def get_countries(self) -> Dict[str, Any]:
        """
        获取支持的国家列表

        Returns:
            国家列表响应
        """
        return self.request('GET', '/api/countries')

    # ========== 订单管理 ==========

    def create_order(self, order_data: Dict[str, Any]) -> Dict[str, Any]:
        """
        创建订单

        Args:
            order_data: 订单数据
                - productId: 产品ID
                - quantity: 数量
                - customerInfo: 客户信息

        Returns:
            订单创建响应
        """
        return self.request('POST', '/api/orders', order_data)

    def get_orders(self, **params) -> Dict[str, Any]:
        """
        获取订单列表

        Args:
            **params: 查询参数
                - status: 订单状态
                - startDate: 开始日期
                - endDate: 结束日期
                - limit: 返回数量限制
                - offset: 偏移量

        Returns:
            订单列表响应
        """
        query_string = urlencode({k: v for k, v in params.items() if v is not None})
        path = f"/api/orders?{query_string}" if query_string else "/api/orders"
        return self.request('GET', path)

    def get_order(self, order_id: str) -> Dict[str, Any]:
        """
        获取订单详情

        Args:
            order_id: 订单ID

        Returns:
            订单详情响应
        """
        return self.request('GET', f"/api/orders/{order_id}")

    # ========== eSIM管理 ==========

    def get_esims(self, **params) -> Dict[str, Any]:
        """
        获取eSIM列表

        Args:
            **params: 查询参数
                - status: eSIM状态
                - orderId: 关联订单ID
                - limit: 返回数量限制
                - offset: 偏移量

        Returns:
            eSIM列表响应
        """
        query_string = urlencode({k: v for k, v in params.items() if v is not None})
        path = f"/api/esims?{query_string}" if query_string else "/api/esims"
        return self.request('GET', path)

    def get_esim(self, esim_id: str) -> Dict[str, Any]:
        """
        获取eSIM详情

        Args:
            esim_id: eSIM ID

        Returns:
            eSIM详情响应
        """
        return self.request('GET', f"/api/esims/{esim_id}")

    def topup_esim(self, esim_id: str, topup_data: Dict[str, Any]) -> Dict[str, Any]:
        """
        eSIM充值

        Args:
            esim_id: eSIM ID
            topup_data: 充值数据
                - packageId: 充值套餐ID
                - amount: 充值金额

        Returns:
            充值响应
        """
        return self.request('POST', f"/api/esims/{esim_id}/topup", topup_data)

    def get_esim_usage(self, esim_id: str) -> Dict[str, Any]:
        """
        获取eSIM使用统计

        Args:
            esim_id: eSIM ID

        Returns:
            使用统计响应
        """
        return self.request('GET', f"/api/esims/{esim_id}/usage")

    # ========== 账户管理 ==========

    def get_account(self) -> Dict[str, Any]:
        """
        获取账户信息

        Returns:
            账户信息响应
        """
        return self.request('GET', '/api/account')

    def get_balance(self) -> Dict[str, Any]:
        """
        获取账户余额

        Returns:
            余额信息响应
        """
        return self.request('GET', '/api/account/balance')

    def get_finance_records(self, **params) -> Dict[str, Any]:
        """
        获取财务记录

        Args:
            **params: 查询参数
                - type: 记录类型（recharge/deduction）
                - startDate: 开始日期
                - endDate: 结束日期
                - limit: 返回数量限制
                - offset: 偏移量

        Returns:
            财务记录响应
        """
        query_string = urlencode({k: v for k, v in params.items() if v is not None})
        path = f"/api/finance-records?{query_string}" if query_string else "/api/finance-records"
        return self.request('GET', path)


def example_usage():
    """使用示例"""
    # 初始化SDK
    sdk = EsimSDK({
        'apiKey': 'your-api-key',
        'apiSecret': 'your-api-secret',
        'baseURL': 'https://api.your-domain.com'
    })

    try:
        # 获取产品列表
        products = sdk.get_products(country='US', limit=10)
        print('产品列表:', json.dumps(products, indent=2, ensure_ascii=False))

        # 创建订单
        order = sdk.create_order({
            'productId': 'product-id',
            'quantity': 1,
            'customerInfo': {
                'name': '张三',
                'email': 'zhangsan@example.com',
                'phone': '+86 138 0000 0000'
            }
        })
        print('订单创建成功:', json.dumps(order, indent=2, ensure_ascii=False))

        # 获取账户信息
        account = sdk.get_account()
        print('账户信息:', json.dumps(account, indent=2, ensure_ascii=False))

    except Exception as e:
        print(f'API调用失败: {e}')


if __name__ == '__main__':
    example_usage()