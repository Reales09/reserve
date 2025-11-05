import 'package:flutter/material.dart';
import 'package:rupu/config/routers/app_bindings.dart';
import 'package:rupu/presentation/views/horizontal_properties/voting_group_detail_view.dart';

class VotingGroupDetailScreen extends StatelessWidget {
  static const name = 'voting-group-detail-screen';
  final int pageIndex;
  final int propertyId;
  final int votingGroupId;

  const VotingGroupDetailScreen({
    super.key,
    required this.pageIndex,
    required this.propertyId,
    required this.votingGroupId,
  });

  @override
  Widget build(BuildContext context) {
    VotingGroupDetailBinding.register(
      propertyId: propertyId,
      votingGroupId: votingGroupId,
    );
    return VotingGroupDetailView(
      propertyId: propertyId,
      votingGroupId: votingGroupId,
    );
  }
}
