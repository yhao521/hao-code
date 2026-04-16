#!/usr/bin/env python3
"""
生成 HC 字母组合风格的代码编辑器图标
蓝色背景 + 白色 HC 字母 + 底部浅蓝装饰线
"""
from pathlib import Path
from PIL import Image, ImageDraw

def create_icon(size, output_path):
    """创建指定尺寸的图标"""
    img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 绘制圆角矩形背景 - 蓝色
    padding = int(size * 0.03)
    corner_radius = int(size * 0.16)
    bg_color = (37, 99, 235)  # #2563eb
    
    draw.rounded_rectangle(
        [padding, padding, size-padding, size-padding],
        radius=corner_radius,
        fill=bg_color
    )
    
    # 计算比例因子
    scale = size / 512.0
    
    # 绘制 H 字母 - 白色实心
    white = (255, 255, 255)
    
    # H 的左竖条
    h_left_bar_x = int(87 * scale)
    h_left_bar_y = int(138 * scale)
    h_left_bar_w = int(28 * scale)
    h_left_bar_h = int(230 * scale)
    h_left_bar_r = int(14 * scale)
    
    draw.rounded_rectangle(
        [h_left_bar_x, h_left_bar_y, h_left_bar_x + h_left_bar_w, h_left_bar_y + h_left_bar_h],
        radius=h_left_bar_r,
        fill=white
    )
    
    # H 的右竖条
    h_right_bar_x = int(217 * scale)
    h_right_bar_y = int(138 * scale)
    h_right_bar_w = int(28 * scale)
    h_right_bar_h = int(230 * scale)
    h_right_bar_r = int(14 * scale)
    
    draw.rounded_rectangle(
        [h_right_bar_x, h_right_bar_y, h_right_bar_x + h_right_bar_w, h_right_bar_y + h_right_bar_h],
        radius=h_right_bar_r,
        fill=white
    )
    
    # H 的中间横条
    h_mid_bar_x = int(87 * scale)
    h_mid_bar_y = int(239 * scale)
    h_mid_bar_w = int(158 * scale)
    h_mid_bar_h = int(28 * scale)
    h_mid_bar_r = int(14 * scale)
    
    draw.rounded_rectangle(
        [h_mid_bar_x, h_mid_bar_y, h_mid_bar_x + h_mid_bar_w, h_mid_bar_y + h_mid_bar_h],
        radius=h_mid_bar_r,
        fill=white
    )
    
    # 绘制 C 字母 - 块状风格（顶部、右侧、底部三个横条）
    # C 的顶部横条
    c_top_x = int(271 * scale)
    c_top_y = int(138 * scale)
    c_top_w = int(154 * scale)
    c_top_h = int(28 * scale)
    c_top_r = int(14 * scale)
    
    draw.rounded_rectangle(
        [c_top_x, c_top_y, c_top_x + c_top_w, c_top_y + c_top_h],
        radius=c_top_r,
        fill=white
    )
    
    # C 的右侧竖条
    c_right_x = int(397 * scale)
    c_right_y = int(138 * scale)
    c_right_w = int(28 * scale)
    c_right_h = int(230 * scale)
    c_right_r = int(14 * scale)
    
    draw.rounded_rectangle(
        [c_right_x, c_right_y, c_right_x + c_right_w, c_right_y + c_right_h],
        radius=c_right_r,
        fill=white
    )
    
    # C 的底部横条
    c_bottom_x = int(271 * scale)
    c_bottom_y = int(340 * scale)
    c_bottom_w = int(154 * scale)
    c_bottom_h = int(28 * scale)
    c_bottom_r = int(14 * scale)
    
    draw.rounded_rectangle(
        [c_bottom_x, c_bottom_y, c_bottom_x + c_bottom_w, c_bottom_y + c_bottom_h],
        radius=c_bottom_r,
        fill=white
    )
    
    # 底部装饰线 - 浅蓝色
    accent_color = (147, 197, 253)  # #93c5fd
    accent_x = int(75 * scale)
    accent_y = int(380 * scale)
    accent_w = int(362 * scale)
    accent_h = int(13 * scale)
    accent_r = int(6 * scale)
    
    draw.rounded_rectangle(
        [accent_x, accent_y, accent_x + accent_w, accent_y + accent_h],
        radius=accent_r,
        fill=accent_color
    )
    
    # 保存图标
    img.save(output_path)
    print(f"✓ 已生成: {output_path} ({size}x{size})")

def main():
    """生成所有尺寸的图标"""
    base_dir = Path(__file__).parent
    darwin_dir = base_dir / "darwin"
    
    # macOS 所需的所有尺寸
    sizes = [16, 32, 64, 128, 256, 512, 1024]
    
    for size in sizes:
        output_path = darwin_dir / f"appicon_{size}.png"
        create_icon(size, output_path)
    
    # 生成主 appicon.png (512x512)
    main_icon = base_dir / "appicon.png"
    create_icon(512, main_icon)
    
    print("\n✓ 所有图标生成完成！")

if __name__ == "__main__":
    main()
