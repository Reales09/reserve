class VotingsResponseModel {
  final bool success;
  final String? message;
  final List<Map<String, dynamic>> data;

  const VotingsResponseModel({
    required this.success,
    this.message,
    required this.data,
  });

  factory VotingsResponseModel.fromJson(Map<String, dynamic> json) {
    return VotingsResponseModel(
      success: json['success'] ?? false,
      message: json['message'],
      data: List<Map<String, dynamic>>.from(json['data'] as List<dynamic>),
    );
  }
}
