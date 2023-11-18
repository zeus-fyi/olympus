import asyncio
import os

from flask import Flask, jsonify, request
from telethon import TelegramClient

app = Flask(__name__)


def get_number():
    return "17575828406"


async def start_telegram_client_async(group_name=None):
    api_id = int(os.environ['TELEGRAM_API_ID'])
    api_hash = os.environ['TELEGRAM_API_HASH']
    client = TelegramClient('session_name', api_id, api_hash)
    client.start(get_number)
    await client.connect()
    messages_data = []
    async for dialog in client.iter_dialogs():
        # Check if the dialog is a group
        if dialog.is_group and str.startswith(dialog.name, group_name):
            print(f'Group name: {dialog.name}')
            messages = await client.get_messages(dialog, limit=10)
            for message in messages:
                messages_data.append({
                    'sender_id': message.sender_id,
                    'message_text': message.text,
                })
                print(f'Message from {message.sender_id}: {message.text}')
    return messages_data


@app.route('/msgs', methods=['GET'])
def start_telegram_client():
    group_name = request.args.get('group_name', default=None)  # Default to None if not provided
    return asyncio.run(start_telegram_client_async(group_name))


@app.route('/health', methods=['GET'])
def health_check():
    return jsonify({'status': 'running'})


if __name__ == '__main__':
    app.run(debug=False, host='0.0.0.0', port=8000)
