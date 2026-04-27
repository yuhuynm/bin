# Frontend Handoff: Secret Share App

## Mục tiêu
Xây dựng frontend cho một app chia sẻ nội dung nhạy cảm qua link.

Người dùng có thể:
- paste nội dung `text`, `env`, hoặc `markdown`
- chọn có password hoặc không
- chọn thời gian hết hạn
- chọn chế độ đọc 1 lần
- tạo link để gửi cho người khác

Backend đã xử lý toàn bộ:
- tạo token
- mã hoá nội dung
- hash password
- kiểm tra hết hạn
- kiểm tra secret đã bị consume chưa
- tăng view count và consume khi one-time

Frontend chỉ cần:
- render form
- gọi API
- hiển thị trạng thái phù hợp

## Luồng sản phẩm

### 1. Trang tạo secret
Màn hình có:
- một textarea lớn để paste nội dung
- select `content type` với các giá trị:
  - `text`
  - `env`
  - `markdown`
- toggle `Require password`
- input password, chỉ hiện khi toggle bật
- select `Expires in`
  - `1 hour`
  - `12 hours`
  - `24 hours`
  - `72 hours`
  - `168 hours`
- checkbox `Burn after first read`
- button `Create secure link`

Khi tạo thành công:
- hiển thị `token`
- hiển thị URL đầy đủ để copy
- nếu có password thì nhắc người dùng gửi password qua kênh khác

### 2. Trang xem secret theo token
Route FE đề xuất:
- `/s/:token`

Luồng:
1. FE đọc `token` từ URL
2. FE gọi `GET /api/v1/secrets/:token`
3. Nếu response trả metadata với `hasPassword = true`
   - hiển thị form nhập password
4. Nếu response trả luôn `content`
   - hiển thị nội dung secret

### 3. Trang unlock secret
Nếu secret có password:
1. người dùng nhập password
2. FE gọi `POST /api/v1/secrets/:token/unlock`
3. nếu đúng password thì hiển thị nội dung secret
4. nếu sai password thì hiện lỗi

## API contract

### 1. Tạo secret
`POST /api/v1/secrets`

Request body:
```json
{
  "content": "DATABASE_URL=postgres://admin:password@localhost:5432/go_bin",
  "contentType": "env",
  "password": "secret123",
  "expiresInHours": 24,
  "oneTime": true
}
```

Lưu ý:
- `password` có thể bỏ trống
- `contentType` hợp lệ: `text`, `env`, `markdown`
- `expiresInHours` tối đa hiện tại là `720`

Success response `201`:
```json
{
  "token": "AbCdEf123456",
  "url": "/api/v1/secrets/AbCdEf123456",
  "hasPassword": true,
  "expiresAt": "2026-04-27T10:00:00Z"
}
```

### 2. Lấy secret theo token
`GET /api/v1/secrets/:token`

Case A: secret không có password, response `200`
```json
{
  "token": "AbCdEf123456",
  "content": "DATABASE_URL=postgres://admin:password@localhost:5432/go_bin",
  "contentType": "env",
  "hasPassword": false,
  "expiresAt": "2026-04-27T10:00:00Z"
}
```

Case B: secret có password, response `200`
```json
{
  "token": "AbCdEf123456",
  "contentType": "markdown",
  "hasPassword": true,
  "requiresAuth": true,
  "expiresAt": "2026-04-27T10:00:00Z",
  "isConsumed": false
}
```

Case C: secret không tồn tại, response `404`
```json
{
  "error": "secret not found"
}
```

Case D: secret hết hạn hoặc đã bị consume, response `410`
```json
{
  "error": "secret unavailable"
}
```

### 3. Unlock secret có password
`POST /api/v1/secrets/:token/unlock`

Request body:
```json
{
  "password": "secret123"
}
```

Success response `200`:
```json
{
  "token": "AbCdEf123456",
  "content": "# Production Notes\n\n- Deploy at 22:00",
  "contentType": "markdown",
  "hasPassword": true,
  "expiresAt": "2026-04-27T10:00:00Z"
}
```

Wrong password `401`:
```json
{
  "error": "invalid password"
}
```

Secret không tồn tại `404`:
```json
{
  "error": "secret not found"
}
```

Secret hết hạn hoặc đã consume `410`:
```json
{
  "error": "secret unavailable"
}
```

## State frontend cần xử lý

### Trang create
State đề xuất:
- `content`
- `contentType`
- `requirePassword`
- `password`
- `expiresInHours`
- `oneTime`
- `isSubmitting`
- `createResult`
- `error`

Validation FE nên có:
- không cho submit nếu `content` rỗng
- nếu bật password thì password không được rỗng
- có thể chặn password quá ngắn dưới 6 ký tự

### Trang view secret
State đề xuất:
- `loading`
- `secretMeta`
- `secretContent`
- `needsPassword`
- `password`
- `unlocking`
- `error`

## Mapping UI theo response

### Khi `GET /api/v1/secrets/:token` trả content
UI hiển thị:
- nội dung secret
- badge `text/env/markdown`
- cảnh báo nếu đây là one-time secret thì có thể đã bị consume ngay sau khi load

### Khi `GET /api/v1/secrets/:token` trả metadata có password
UI hiển thị:
- input password
- button `Unlock secret`

### Khi nhận `404`
UI hiển thị:
- `Secret not found`

### Khi nhận `410`
UI hiển thị:
- `This secret has expired or was already consumed`

### Khi nhận `401`
UI hiển thị:
- `Wrong password`

## Rendering nội dung

### `contentType = text`
- render trong `pre` hoặc `textarea readonly`

### `contentType = env`
- render monospace
- giữ nguyên line breaks
- có nút `Copy`

### `contentType = markdown`
- v1 có thể render plain text trước
- nếu render markdown HTML thì phải sanitize

Khuyến nghị cho v1:
- render tất cả dưới dạng plain text trong khối `pre-wrap`
- chưa cần markdown preview ngay

## Các lưu ý quan trọng cho FE
- không lưu secret content vào localStorage
- không log secret content ra console trong production
- không đưa secret content vào query string
- sau khi copy xong có thể hiện toast `Copied`
- nên có nút `Create another secret`
- nên phân biệt rõ trạng thái `loading`, `empty`, `error`, `success`

## Gợi ý UX
- form create nên nằm bên trái, panel kết quả nằm bên phải trên desktop
- trên mobile thì stack theo chiều dọc
- textarea phải đủ lớn để paste `.env`
- password input nên có nút show/hide
- sau khi tạo xong nên auto focus vào ô chứa URL hoặc có nút copy rõ ràng

## Prompt mẫu để giao cho người làm FE
Bạn có thể copy prompt này:

```text
Hãy build frontend cho app chia sẻ secret bằng link.

Tech stack frontend:
- tự chọn stack hiện đại phù hợp
- UI sạch, rõ, ưu tiên dễ dùng

Yêu cầu màn hình:
1. Trang create secret
- textarea để paste nội dung
- select content type: text, env, markdown
- toggle require password
- nếu bật thì hiện input password
- select expire time
- checkbox burn after first read
- button create secure link
- sau khi tạo thành công hiển thị URL để copy và thông tin hasPassword, expiresAt

2. Trang view secret tại route /s/:token
- khi vào trang, gọi GET /api/v1/secrets/:token
- nếu response có content thì hiển thị secret luôn
- nếu response có hasPassword=true và requiresAuth=true thì hiển thị form nhập password
- submit password bằng POST /api/v1/secrets/:token/unlock
- nếu unlock thành công thì hiển thị nội dung

API contract:
- POST /api/v1/secrets
- GET /api/v1/secrets/:token
- POST /api/v1/secrets/:token/unlock

Request tạo secret:
{
  "content": "DATABASE_URL=postgres://admin:password@localhost:5432/go_bin",
  "contentType": "env",
  "password": "secret123",
  "expiresInHours": 24,
  "oneTime": true
}

Response create:
{
  "token": "AbCdEf123456",
  "url": "/api/v1/secrets/AbCdEf123456",
  "hasPassword": true,
  "expiresAt": "2026-04-27T10:00:00Z"
}

Lưu ý UX và security:
- không lưu secret vào localStorage
- render plain text an toàn cho text/env/markdown ở v1
- có loading, error, success states rõ ràng
- có nút copy
- tối ưu tốt cho mobile và desktop
```

## Scope v1 nên giữ gọn
- chưa cần auth user
- chưa cần dashboard quản lý list secrets
- chưa cần markdown preview nâng cao
- chưa cần analytics
- chưa cần client-side encryption

Chỉ cần làm tốt:
- create secret
- view secret
- unlock secret
- copy link
- hiển thị trạng thái lỗi rõ ràng
