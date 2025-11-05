import 'package:rupu/domain/entities/voting_option.dart';

class VotingOptionsResult {
  final bool success;
  final String? message;
  final List<VotingOption> data;

  const VotingOptionsResult({
    required this.success,
    this.message,
    required this.data,
  });
}
