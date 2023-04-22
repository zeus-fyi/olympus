import tiktoken
from flask import Flask, request, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app, resources={r"*": {"origins": "http://localhost:3000"}})


@app.route('/token-count', methods=['POST'])
def token_count():
    try:
        text = request.get_json().get('text', '')
        encoding = tiktoken.encoding_for_model('gpt2')
        count = len(encoding.encode(text))
        return jsonify({'count': count}), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500


if __name__ == '__main__':
    app.run(debug=False, host='0.0.0.0', port=5001)
