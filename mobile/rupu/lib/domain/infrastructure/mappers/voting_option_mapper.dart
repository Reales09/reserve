import 'package:rupu/domain/entities/voting_option.dart';
import 'package:rupu/domain/entities/voting_options_result.dart';

class VotingOptionMapper {
  static VotingOptionsResult fromJson(Map<String, dynamic> json) {
    return VotingOptionsResult(
      success: json['success'] ?? false,
      message: json['message'],
      data: List<VotingOption>.from(
        (json['data'] as List<dynamic>).map(
          (x) => VotingOption(
            id: x['id'],
            votingId: x['voting_id'],
            optionText: x['option_text'],
            optionCode: x['option_code'],
            color: x['color'],
            displayOrder: x['display_order'],
            isActive: x['is_active'],
          ),
        ),
      ),
    );
  }
}
