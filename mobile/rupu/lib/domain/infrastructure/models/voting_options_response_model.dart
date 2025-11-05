class VotingOptionsResponseModel {
  final bool success;
  final String? message;
  final List<Map<String, dynamic>> data;

  const VotingOptionsResponseModel({
    required this.success,
    this.message,
    required this.data,
  });

  factory VotingOptionsResponseModel.fromJson(Map<String, dynamic> json) {
    return VotingOptionsResponseModel(
      success: json['success'] ?? false,
      message: json['message'],
      data: List<Map<String, dynamic>>.from(json['data'] as List<dynamic>),
    );
  }
}
